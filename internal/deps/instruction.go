package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
)

type insSet map[*instruction]struct{}

type instruction struct {
	DynAddress model.Address
	Instr      repr.Instruction

	trueDepsFwd    insSet
	trueDepsBack   insSet
	antiDepsFwd    insSet
	antiDepsBack   insSet
	outputDepsFwd  insSet
	outputDepsBack insSet

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

		blockIdx: index,
	}
}

func (i *instruction) ptr() Instruction { return Instruction{i: i} }

// Idx returns index of an instruction in its basic block.
func (i *instruction) Idx() int { return i.blockIdx }
