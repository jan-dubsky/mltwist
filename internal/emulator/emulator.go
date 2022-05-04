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
	prog      *deps.Program
	ip        model.Addr
	stateProv StateProvider

	state *state.State
}

// New creates new emulator instance
func New(prog *deps.Program, ip model.Addr, stateProv StateProvider) *Emulator {
	return &Emulator{
		prog:      prog,
		ip:        ip,
		stateProv: stateProv,

		state: state.New(),
	}
}

// IP returns current state of instruction pointer of the program.
func (e *Emulator) IP() model.Addr { return e.ip }

// State returns the internal instate of emulator.
//
// The state doesn't have to be treated as readonly, but its modifications will
// be reflected in further steps of the emulation. analogously, any further step
// of emulation will modify the state returned. Consequently, all reads and
// modifications of the returned state must be synchronized with Step calls.
func (e *Emulator) State() *state.State { return e.state }

// Step performs a single instruction step of an emulation.
func (e *Emulator) Step(ins Instruction) Evaluation {
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
	return eval
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
	val = exprtransform.SetWidthConst(val, w)

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

	for _, intv := range e.state.Mems.Missing(key, addr, w) {
		addr, w := intv.Begin(), intv.Width()
		val := e.stateProv.Memory(key, intv.Begin(), w)
		val = exprtransform.SetWidthConst(val, w)

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
