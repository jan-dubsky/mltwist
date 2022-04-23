package deps

import (
	"mltwist/internal/repr"
	"mltwist/pkg/model"
)

type insSet map[*instruction]struct{}

type instruction struct {
	DynAddress model.Addr
	Instr      repr.Instruction

	trueDepsFwd     insSet
	trueDepsBack    insSet
	antiDepsFwd     insSet
	antiDepsBack    insSet
	outputDepsFwd   insSet
	outputDepsBack  insSet
	controlDepsFwd  insSet
	controlDepsBack insSet
	specialDepsFwd  insSet
	specialDepsBack insSet

	blockIdx int
}

func newInstruction(ins repr.Instruction, index int) *instruction {
	// Those are absolutely thumbsucked numbers of expected dependencies.
	// There is no scientific neither measured reason for those constant,
	// but given that Go default for map size is 100, those are definitely
	// more optimal sizes for dependency maps.
	const expectedDeps = 2

	return &instruction{
		DynAddress: ins.Address,
		Instr:      ins,

		trueDepsFwd:    make(insSet, expectedDeps),
		trueDepsBack:   make(insSet, expectedDeps),
		antiDepsFwd:    make(insSet, expectedDeps),
		antiDepsBack:   make(insSet, expectedDeps),
		outputDepsFwd:  make(insSet, expectedDeps),
		outputDepsBack: make(insSet, expectedDeps),
		// Each non-jump instruction is dependent on exactly one
		// instruction in forward direction - the jump/call instruction
		// at the end of basic block.
		controlDepsFwd: make(insSet, 1),
		// Here optimal size is the size of basic block which is not
		// known to us for jump/call instruction at the end of the block
		// and 0 for any non-jump instruction. As there is significantly
		// more more non-jump than jump instructions, we use 0 not to
		// waste memory and we rely on exponential size increasing in
		// case of jump/call instructions.
		controlDepsBack: make(insSet, 0),

		// No special instructions are expected as those are very rare.
		specialDepsFwd:  make(insSet, 0),
		specialDepsBack: make(insSet, 0),

		blockIdx: index,
	}
}

// Idx returns index of an instruction in its basic block.
func (i *instruction) Idx() int { return i.blockIdx }

func (i *instruction) setIndex(idx int) { i.blockIdx = idx }
