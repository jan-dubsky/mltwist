package deps

import (
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type insSet map[*instruction]struct{}

type regSet map[expr.Key]struct{}

type instruction struct {
	// typ is an instruction special type.
	typ model.Type
	// origAddr is original address of the instruction in the code. This
	// value remains constant even in case the instruction is moved to other
	// place.
	origAddr model.Addr

	// bytes are instruction bytes in the code.
	bytes []byte
	// details provide platform-specific instruction methods.
	details model.PlatformDetails

	// effects contains all side effects the instruction has.
	effects []expr.Effect
	// jumpTargets is list of expressions which describe jump targets of
	// this instruction. If this array is non-empty, the instruction is some
	// form of jump (control flow instruction).
	//
	// Constant expressions jumping to following instructions are omitted as
	// those are not real jumps. This is important as those expressions will
	// be typically constants which cannot be adjusted by instruction
	// moving. So if this array contained jumps to the following
	// instructions as well, those could be interpreted as real jump if the
	// instruction would be moved to other position (memory address).
	jumpTargets []expr.Expr

	// currAddr is current address of the instruction in the moved code.
	currAddr model.Addr

	// inRegs is a set of instruction input (loaded) registers.
	inRegs regSet
	// outRegs is a set of registers written by the instruction.
	outRegs regSet

	// loads is list of all memory loads the instruction does in all its
	// effects.
	loads []expr.MemLoad
	// stores is list of all memory stores the instruction does.
	stores []expr.MemStore

	// depsFwd is a set of references to all instructions in the basic block
	// which have to be executed after this instruction.
	depsFwd insSet
	// depsBack is a set of references to all instructions in the basic
	// block which have to be executed before this instruction.
	depsBack insSet

	// blockIdx is index of instruction in a basic block.
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

		currAddr: ins.Addr,

		inRegs:  inputRegs(ins.Effects),
		outRegs: outputRegs(ins.Effects),

		loads:  loads(ins.Effects),
		stores: stores(ins.Effects),

		depsFwd:  make(insSet, 5),
		depsBack: make(insSet, 5),

		blockIdx: index,
	}
}

// Idx returns index of an instruction in its basic block.
func (i *instruction) Idx() int { return i.blockIdx }

// Begin returns the in-memory address of the instruction in the current order
// of the program - i.e. after all instruction moves.
func (i *instruction) Begin() model.Addr { return i.currAddr }

// Len returns length of an instruction in bytes.
func (i *instruction) Len() model.Addr { return model.Addr(len(i.bytes)) }

// End returns memory address of an instruction following this instruction. The
// address taken into account is the one returned by Begin(), not OrigAddr().
func (i *instruction) End() model.Addr { return i.currAddr + i.Len() }

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
			regs[l.Key()] = struct{}{}
		}
	}

	return regs
}

func outputRegs(effects []expr.Effect) regSet {
	regs := make(regSet, 1)

	for _, effect := range effects {
		if e, ok := effect.(expr.RegStore); ok {
			regs[e.Key()] = struct{}{}
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
