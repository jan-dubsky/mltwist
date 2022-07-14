package emulator

import (
	"fmt"
	"mltwist/internal/deps"
	"mltwist/internal/exprtransform"
	"mltwist/internal/state"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Emulator struct {
	code      *deps.Code
	stateProv StateProvider

	// State represents the current state of the emulation.
	//
	// It's allowed to both read and write values from this state. The only
	// requirement on values written is that all those values have to be of
	// type expr.Const. Violation of this requirement might result in
	// undefined behaviour of the emulation.
	//
	// As this field is used by emulation, it's required that any access of
	// State must be synchronized with Step method calls. Any violation of
	// this synchronization might result in undefined behaviour.
	State *state.State
}

// New creates new emulator instance.
//
// The newly created emulator emulates instructions in prog and starts emulation
// at address ip. The initial state of program memory and registers is given by
// state. If value of any memory bytes or a register is unknown to the emulator,
// such value is obtained using stateProv.
//
// Argument prog is treated as read-only value, but it must not be modified
// during the Emulator lifetime.
func New(
	code *deps.Code,
	ip model.Addr,
	stateProv StateProvider,
	state *state.State,
) *Emulator {
	state.Regs.Store(expr.IPKey, expr.ConstFromUint(ip), model.AddrWidth)

	return &Emulator{
		code:      code,
		stateProv: stateProv,
		State:     state,
	}
}

// MustIP returns a current value of instruction pointer.
//
// This function panics if the value of instruction pointer is not stored in a
// state, if the value stored there is not constant or if the constant doesn't
// fit model.Addr type.
func (e *Emulator) MustIP() model.Addr {
	ex, ok := e.State.Regs.Load(expr.IPKey, model.AddrWidth)
	if !ok {
		panic("instruction pointer register is missing in registry file.")
	}

	c := ex.(expr.Const)
	ip, ok := expr.ConstUint[model.Addr](c)
	if !ok {
		panic(fmt.Sprintf(
			"constant value of instruction pointer doesn't fit address: %v",
			c.Bytes(),
		))
	}

	return ip
}

// instruction returns instruction currently pointer by the instruction pointer.
func (e *Emulator) instruction(ip model.Addr) (deps.Instruction, error) {
	block, ok := e.code.Address(ip)
	if !ok {
		err := fmt.Errorf("cannot find block at address 0x%x", ip)
		return deps.Instruction{}, err
	}

	ins, ok := block.Address(ip)
	if !ok {
		err := fmt.Errorf("cannot find instruction at address 0x%x", ip)
		return deps.Instruction{}, err
	}

	return ins, nil
}

// Step performs a single instruction step of an emulation.
func (e *Emulator) Step() (*Step, error) {
	ip := e.MustIP()
	ins, err := e.instruction(ip)
	if err != nil {
		return nil, err
	}

	efs := ins.Effects()
	s := newStep(len(efs))

	efs = exprtransform.EffectsApply(efs, func(ex expr.Expr) expr.Expr {
		return e.eval(ex, s)
	})

	var jumped bool
	for _, ef := range efs {
		// We can calculate new value of instruction pointer from our
		// effects directly. The reason is that jump instruction can be
		// never moved from it's last position n=in a basic block due to
		// control dependencies. As we don't allow basic blocks to be
		// moved in the address space, we are certain that all addresses
		// (including constants pointing to the following instruction)
		// are still valid.
		if rStore, ok := ef.(expr.RegStore); ok && rStore.Key() == expr.IPKey {
			jumped = true
		}

		s.recordOutput(ef)
		ok := e.State.Apply(ef)
		if !ok {
			panic(fmt.Errorf("bug: state change couldn't be applied: %v", ef))
		}
	}

	// Instruction is not a jump instruction so we have to adjust
	// instruction pointer ourselves.
	if !jumped {
		c := expr.ConstFromUint(ins.End())
		e.State.Regs.Store(expr.IPKey, c, model.AddrWidth)
	}

	return s, nil
}

// eval substitutes all non-constant expressions (register loads and memory
// loads) for values from the emulator storage. If the value read is not stored
// in the storage, the value is supplied by stateProvider interface and stored
// into the storage.
func (e *Emulator) eval(ex expr.Expr, s *Step) expr.Const {
	ex = e.evalRegsFully(ex, s)
	ex = e.evalMemoryFully(ex, s)

	// We have replaced all register and memory loads by constants,
	// so the result of const fold has to be constant.
	return exprtransform.ConstFold(ex).(expr.Const)
}

// regValue reads constant value of register key of width w from the internal
// storage. Is the value is not present in the register storage, it's obtained
// using stateProvider interface and stored in the register storage.
func (e *Emulator) regValue(key expr.Key, w expr.Width) expr.Const {
	if val, ok := e.State.Regs.Load(key, w); ok {
		return val.(expr.Const)
	}

	val := e.stateProv.Register(key, w)
	val = val.WithWidth(w)
	e.State.Regs.Store(key, val, w)

	return val
}

func (e *Emulator) evalRegsFully(ex expr.Expr, s *Step) expr.Expr {
	ex = exprtransform.ReplaceAll(ex, func(curr expr.RegLoad) (expr.Expr, bool) {
		key := curr.Key()
		val := e.regValue(key, curr.Width())

		s.inputReg(key, val)
		return val, true
	})

	// Constant folding makes no sense here as we will constant-fold all
	// addresses in memory evaluation and we will as well evaluate the final
	// result after memory evaluation.
	return ex
}

// memValue reads a constant value of memory address adds in memory address
// space key of width w from the internal storage. Is the value is not present
// in the memory, it's obtained using stateProvider interface and stored in the
// memory storage.
func (e *Emulator) memValue(key expr.Key, addr model.Addr, w expr.Width) expr.Const {
	if val, ok := e.State.Mems.Load(key, addr, w); ok {
		return val.(expr.Const)
	}

	for _, intv := range e.State.Mems.Missing(key, addr, w).Intervals() {
		addr := intv.Begin()
		// We are certain that length first expr.Width as we provided
		// expr.Width to the missing call. The wort possible case is
		// that the whole interval is missing, but the interval missing
		// is then w long at most and w is expr.Width.
		w := expr.Width(intv.Len())

		val := e.stateProv.Memory(key, intv.Begin(), w)
		val = val.WithWidth(w)

		e.State.Mems.Store(key, addr, val, w)
	}

	val, ok := e.State.Mems.Load(key, addr, w)
	if !ok {
		panic(fmt.Sprintf(
			"bug: memory with width %d at addr 0x%x not present",
			w, addr))
	}

	return val.(expr.Const)
}

func (e *Emulator) evalMemoryFully(ex expr.Expr, s *Step) expr.Expr {
	ex = exprtransform.ReplaceAll(ex, func(curr expr.MemLoad) (expr.Expr, bool) {
		key, w := curr.Key(), curr.Width()

		// We don't check for ok as this cast must be always successful.
		// Remember, we are evaluating the most bottom expression and
		// all registers are already evaluated -> this MUST be constant.
		addrConst := exprtransform.ConstFold(curr.Addr()).(expr.Const)
		addr, _ := expr.ConstUint[model.Addr](addrConst)

		val := e.memValue(key, addr, w)
		s.memRead(key, addr, val)

		return val, true
	})

	return ex
}
