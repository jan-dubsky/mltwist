package emulator

import (
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// RegSet if a set of registers and their respective values loaded or stored.
type RegSet map[expr.Key]expr.Const

// MemAccess describes a single memory access the simulation performed.
type MemAccess struct {
	// Key is a key of memory address space accessed.
	Key expr.Key
	// Addr is the address in the memory accessed.
	Addr model.Addr
	// Value is the value loader or stored into the memory.
	Value expr.Const
}

func newMemAccess(key expr.Key, addr model.Addr, c expr.Const) MemAccess {
	return MemAccess{
		Key:   key,
		Addr:  addr,
		Value: c,
	}
}

// Width returns width of the memory access.
func (a MemAccess) Width() expr.Width { return a.Value.Width() }

// Step describes all loads and stores a single emulated instruction performed.
type Step struct {
	// RegLoads contains a set of register keys loaded and their respective
	// values.
	RegLoads RegSet
	// RegStores contains a set of register keys stored and their respective
	// values.
	RegStores RegSet

	// MemLoads contains a list of memory loads.
	MemLoads []MemAccess
	// MemStores contains a list of memory stores.
	MemStores []MemAccess
}

func newStep(numEffects int) *Step {
	return &Step{
		RegLoads:  make(RegSet, 16),
		RegStores: make(RegSet, numEffects),
	}
}

func (e *Step) inputReg(key expr.Key, c expr.Const) { e.RegLoads[key] = c }
func (e *Step) memRead(key expr.Key, addr model.Addr, c expr.Const) {
	e.MemLoads = append(e.MemLoads, newMemAccess(key, addr, c))
}

func (e *Step) recordOutput(effect expr.Effect) {
	switch ef := effect.(type) {
	case expr.MemStore:
		addr, _ := expr.ConstUint[model.Addr](ef.Addr().(expr.Const))
		val := ef.Value().(expr.Const).WithWidth(ef.Width())
		e.MemStores = append(e.MemStores, newMemAccess(ef.Key(), addr, val))
	case expr.RegStore:
		e.RegStores[ef.Key()] = ef.Value().(expr.Const)
	}
}
