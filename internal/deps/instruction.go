package deps

import (
	"mltwist/internal/exprtransform"
	"mltwist/internal/parser"
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

func newInstruction(ins parser.Instruction) *instruction {
	return &instruction{
		typ:      ins.Type,
		origAddr: ins.Addr,
		bytes:    ins.Bytes,
		details:  ins.Details,

		effects:     ins.Effects,
		jumpTargets: jumps(ins),

		currAddr: ins.Addr,

		inRegs:  inputRegs(ins.Effects),
		outRegs: outputRegs(ins.Effects),

		loads:  loads(ins.Effects),
		stores: stores(ins.Effects),

		depsFwd:  make(insSet, 5),
		depsBack: make(insSet, 5),

		blockIdx: -1,
	}
}

// Jumps extracts all expressions the instruction can jump to. Jumps to address
// following the instruction (to address of End()) are filtered away as those
// are not read jump addresses.
func jumps(ins parser.Instruction) []expr.Expr {
	var jumpAddrs []expr.Expr
	for _, ef := range ins.Effects {
		e, ok := ef.(expr.RegStore)
		if !ok {
			continue
		}

		if e.Key() != expr.IPKey {
			continue
		}

		addrs := exprtransform.Possibilities(e.Value())

		// Filter those jump addresses which jump to the following
		// instruction as those are technically not jumps.
		j := 0
		for i := 0; i < len(addrs); i, j = i+1, j+1 {
			a := exprtransform.ConstFold(addrs[i])
			addrs[j] = a

			c, ok := a.(expr.Const)
			if !ok {
				continue
			}

			addr, _ := expr.ConstUint[model.Addr](c)
			if addr != ins.End() {
				continue
			}

			j--
		}
		jumpAddrs = append(jumpAddrs, addrs[:j]...)
	}

	return jumpAddrs
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

// Effects returns a read-only list of all side effects if an instruction.
func (i *instruction) Effects() []expr.Effect { return i.effects }

// Jumps returns a read-only list of all memory addresses this instruction can
// jump to.
func (i *instruction) Jumps() []expr.Expr { return i.jumpTargets }

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
