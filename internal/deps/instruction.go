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
	typ      model.Type
	origAddr model.Addr
	bytes    []byte
	details  model.PlatformDetails

	effects     []expr.Effect
	jumpTargets []expr.Expr

	currAddr model.Addr

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
		typ:      ins.Type,
		origAddr: ins.Addr,
		bytes:    ins.Bytes,
		details:  ins.Details,

		effects:     ins.Effects,
		jumpTargets: ins.JumpTargets,

		inRegs:  inputRegs(ins.Effects),
		outRegs: outputRegs(ins.Effects),

		loads:  loads(ins.Effects),
		stores: stores(ins.Effects),

		depsFwd:  make(insSet, 5),
		depsBack: make(insSet, 5),

		blockIdx: index,
		currAddr: ins.Addr,
	}
}

// Idx returns index of an instruction in its basic block.
func (i *instruction) Idx() int { return i.blockIdx }

// Len returns length of an instruction in bytes.
func (i *instruction) Len() model.Addr { return model.Addr(len(i.bytes)) }

// NextAddr returns memory address of an instruction following this instruction.
// The address taken into account is the one returned by Addr(), not OrigAddr().
func (i *instruction) NextAddr() model.Addr { return i.currAddr + i.Len() }

// Addr returns the in-memory address of the instruction in the current order of
// the program - i.e. after all instruction moves.
func (i *instruction) Addr() model.Addr { return i.currAddr }

// OrigAddr returns the in-memory address of the instruction in the original
// binary.
func (i *instruction) OrigAddr() model.Addr { return i.origAddr }

// Effects returns a list of all side effects if an instruction.
func (i *instruction) Effects() []expr.Effect { return i.effects }

func (i *instruction) setIndex(idx int)     { i.blockIdx = idx }
func (i *instruction) setAddr(a model.Addr) { i.currAddr = a }

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

func addDep(first, second *instruction) {
	first.depsFwd[second] = struct{}{}
	second.depsBack[first] = struct{}{}
}
