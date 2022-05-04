package emulator

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// RegSet if a set of registers and their respective values.
type RegSet map[expr.Key]expr.Const

type MemAccess struct {
	Key   expr.Key
	Addr  model.Addr
	Value expr.Const
}

func newMemAccess(key expr.Key, addr model.Addr, c expr.Const) MemAccess {
	return MemAccess{
		Key:   key,
		Addr:  addr,
		Value: c,
	}
}

func (a MemAccess) Width() expr.Width { return a.Value.Width() }

type Evaluation struct {
	InputRegs  RegSet
	OutputRegs RegSet

	InMem  []MemAccess
	OutMem []MemAccess
}

func newEvaluation() *Evaluation {
	return &Evaluation{
		InputRegs:  make(RegSet, 16),
		OutputRegs: make(RegSet, 4),
	}
}

func (e *Evaluation) inputReg(key expr.Key, c expr.Const) { e.InputRegs[key] = c }
func (e *Evaluation) memRead(key expr.Key, addr model.Addr, c expr.Const) {
	e.InMem = append(e.InMem, newMemAccess(key, addr, c))
}

func (e *Evaluation) recordOutput(effect expr.Effect) {
	switch ef := effect.(type) {
	case expr.MemStore:
		addr, _ := expr.ConstUint[model.Addr](ef.Addr().(expr.Const))
		val := exprtransform.SetWidthConst(ef.Value().(expr.Const), ef.Width())
		e.OutMem = append(e.OutMem, newMemAccess(ef.Key(), addr, val))
	case expr.RegStore:
		e.OutputRegs[ef.Key()] = ef.Value().(expr.Const)
	}
}
