package deps

import (
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type insSet map[*instruction]struct{}

type regSet map[string]struct{}

type instruction struct {
	DynAddress model.Addr
	Instr      basicblock.Instruction

	inRegs  regSet
	outRegs regSet

	loads  []expr.MemLoad
	stores []expr.MemStore

	depsFwd  insSet
	depsBack insSet

	blockIdx int
}

func newInstruction(ins basicblock.Instruction, index int) *instruction {
	return &instruction{
		DynAddress: ins.Addr,
		Instr:      ins,

		inRegs:  inputRegs(ins.Effs),
		outRegs: outputRegs(ins.Effs),

		loads:  loads(ins.Effs),
		stores: stores(ins.Effs),

		depsFwd:  make(insSet, 5),
		depsBack: make(insSet, 5),

		blockIdx: index,
	}
}

// Idx returns index of an instruction in its basic block.
func (i *instruction) Idx() int { return i.blockIdx }

func (i *instruction) setIndex(idx int) { i.blockIdx = idx }

func inputRegs(effects []expr.Effect) regSet {
	// The 2 default value might be too little, but it's reasonable
	// thumbsuck - even CISC architectures typically use at most 3
	// registers. If we omitted the constant, the map would be 100 elements
	// big. This is a better option.
	regs := make(regSet, 2)

	for _, ex := range exprtransform.ExprsMany(effects) {
		for _, l := range exprtransform.FindAll[expr.RegLoad](ex) {
			regs[string(l.Key())] = struct{}{}
		}
	}

	return regs
}

func outputRegs(effects []expr.Effect) regSet {
	regs := make(regSet, 1)

	for _, effect := range effects {
		if e, ok := effect.(expr.RegStore); ok {
			regs[string(e.Key())] = struct{}{}
		}
	}

	return regs
}

func loads(effects []expr.Effect) []expr.MemLoad {
	var loads []expr.MemLoad
	for _, ex := range exprtransform.ExprsMany(effects) {
		loads = append(loads, exprtransform.FindAll[expr.MemLoad](ex)...)
	}
	return loads
}

func stores(effects []expr.Effect) []expr.MemStore {
	var stores []expr.MemStore
	for _, ef := range effects {
		if e, ok := ef.(expr.MemStore); ok {
			stores = append(stores, e)
		}
	}
	return stores
}
