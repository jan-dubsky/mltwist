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
	code  *deps.Code
	ip    model.Addr
	state *state.State

	stateProv StateProvider
}

// New creates new emulator instance.
//
// The newly created emulator emulates instructions in prog and starts emulation
// at address ip. The initial state of program memory and registers is given by
// stat. If value of any memory bytes or a register is unknown to the emulator,
// such value is obtained using stateProv.
//
// Argument prog is treated as read-only value, but it must not be modified
// during the Emulator lifetime.
//
// Argument stat will be used by the Emulator to store writes the emulated
// program performs. Consequently it can be used to read/change the state of the
// emulation. For the very same reason, any access of stat must be synchronized
// with Step method calls. Any violation of this synchronization might result in
// undefined behaviour.
func New(
	code *deps.Code,
	ip model.Addr,
	stat *state.State,
	stateProv StateProvider,
) *Emulator {
	return &Emulator{
		code:      code,
		ip:        ip,
		stateProv: stateProv,
		state:     stat,
	}
}

// IP returns current state of instruction pointer of the program.
func (e *Emulator) IP() model.Addr { return e.ip }

// instruction returns instruction currently pointer by the instruction pointer.
func (e *Emulator) instruction() (deps.Instruction, error) {
	block, ok := e.code.Address(e.ip)
	if !ok {
		err := fmt.Errorf("cannot find block at address 0x%x", e.ip)
		return deps.Instruction{}, err
	}

	ins, ok := block.Address(e.ip)
	if !ok {
		err := fmt.Errorf("cannot find instruction at address 0x%x", e.ip)
		return deps.Instruction{}, err
	}

	return ins, nil
}

// Step performs a single instruction step of an emulation.
func (e *Emulator) Step() (*Step, error) {
	ins, err := e.instruction()
	if err != nil {
		return nil, err
	}

	efs := ins.Effects()
	s := newStep(len(efs))

	efs = exprtransform.EffectsApply(efs, func(ex expr.Expr) expr.Expr {
		return e.eval(ex, s)
	})

	nextIP := e.ip + ins.Len()
	for _, ef := range efs {
		// We can calculate new value of instruction pointer from our
		// effects directly. The reason is that jump instruction can be
		// never moved from it's last position n=in a basic block due to
		// control dependencies. As we don't allow basic blocks to be
		// moved in the address space, we are certain that all addresses
		// (including constants pointing to the following instruction)
		// are still valid.
		if rStore, ok := ef.(expr.RegStore); ok && rStore.Key() == expr.IPKey {
			val := rStore.Value().(expr.Const)
			nextIP, _ = expr.ConstUint[model.Addr](val)

			// We don't want the instruction pointer to be present
			// in register storage as such a value would be ignored
			// by future steps of the emulation.
			continue
		}

		s.recordOutput(ef)
		ok := e.state.Apply(ef)
		if !ok {
			panic(fmt.Errorf("bug: state change couldn't be applied: %v", ef))
		}
	}

	e.ip = nextIP
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
	if val, ok := e.state.Regs.Load(key, w); ok {
		return val.(expr.Const)
	}

	val := e.stateProv.Register(key, w)
	val = val.WithWidth(w)
	e.state.Regs.Store(key, val, w)

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
	if val, ok := e.state.Mems.Load(key, addr, w); ok {
		return val.(expr.Const)
	}

	for _, intv := range e.state.Mems.Missing(key, addr, w).Intervals() {
		addr := intv.Begin()
		// We are certain that length first expr.Width as we provided
		// expr.Width to the missing call. The wort possible case is
		// that the whole interval is missing, but the interval missing
		// is then w long at most and w is expr.Width.
		w := expr.Width(intv.Len())

		val := e.stateProv.Memory(key, intv.Begin(), w)
		val = val.WithWidth(w)

		e.state.Mems.Store(key, addr, val, w)
	}

	val, ok := e.state.Mems.Load(key, addr, w)
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
