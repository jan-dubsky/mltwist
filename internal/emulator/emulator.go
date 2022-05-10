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
	prog  *deps.Program
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
// emulation. For the very same reason, an yaccess of stat must be synchronized
// with Step method calls.
func New(
	prog *deps.Program,
	ip model.Addr,
	stat *state.State,
	stateProv StateProvider,
) *Emulator {
	return &Emulator{
		prog:      prog,
		ip:        ip,
		stateProv: stateProv,
		state:     stat,
	}
}

// IP returns current state of instruction pointer of the program.
func (e *Emulator) IP() model.Addr { return e.ip }

// Step performs a single instruction step of an emulation.
func (e *Emulator) Step() (Evaluation, error) {
	ins, ok := e.prog.AddressIns(e.ip)
	if !ok {
		err := fmt.Errorf("cannot find instruction at address 0x%x", e.ip)
		return Evaluation{}, err
	}

	efs := ins.Effects()

	eval := Evaluation{
		InputRegs:  make(RegSet, 16),
		OutputRegs: make(RegSet, len(efs)),
	}

	efs = exprtransform.EffectsApply(efs, func(ex expr.Expr) expr.Expr {
		return e.eval(ex, &eval)
	})

	nextIP := e.ip + ins.Len()
	for _, ef := range efs {
		if rStore, ok := ef.(expr.RegStore); ok && rStore.Key() == expr.IPKey {
			val := rStore.Value().(expr.Const)
			nextIP, _ = expr.ConstUint[model.Addr](val)
		}

		eval.recordOutput(ef)
		ok := e.state.Apply(ef)
		if !ok {
			panic(fmt.Errorf("bug: state change couldn't be applied: %v", ef))
		}
	}

	e.ip = nextIP
	return eval, nil
}

func (e *Emulator) eval(ex expr.Expr, eval *Evaluation) expr.Const {
	ex = e.evalRegsFully(ex, eval)
	ex = e.evalMemoryFully(ex, eval)

	// We have replaced all register and memory loads by constants,
	// so the result of const fold has to be constant.
	return exprtransform.ConstFold(ex).(expr.Const)
}

func (e *Emulator) regValue(key expr.Key, w expr.Width) expr.Const {
	if val, ok := e.state.Regs[key]; ok {
		return val.(expr.Const)
	}

	val := e.stateProv.Register(key, w)
	val = val.WithWidth(w)

	e.state.Regs[key] = val
	return val
}

func (e *Emulator) evalRegsFully(ex expr.Expr, eval *Evaluation) expr.Expr {
	ex = exprtransform.ReplaceAll(ex, func(curr expr.RegLoad) (expr.Expr, bool) {
		key := curr.Key()
		val := e.regValue(key, curr.Width())

		eval.inputReg(key, val)
		return val, true
	})

	// Constant folding makes no sense here as we will constant-fold all
	// addresses in memory evaluation and we will as well evaluate the final
	// result after memory evaluation.
	return ex
}

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

func (e *Emulator) evalMemoryFully(ex expr.Expr, eval *Evaluation) expr.Expr {
	ex = exprtransform.ReplaceAll(ex, func(curr expr.MemLoad) (expr.Expr, bool) {
		key, w := curr.Key(), curr.Width()

		// We don't check for ok as this cast must be always successful.
		// Remember, we are evaluating the most bottom expression and
		// all registers are already evaluated -> this MUST be constant.
		addrConst := exprtransform.ConstFold(curr.Addr()).(expr.Const)
		addr, _ := expr.ConstUint[model.Addr](addrConst)

		val := e.memValue(key, addr, w)
		eval.memRead(key, addr, val)

		return val, true
	})

	return ex
}
