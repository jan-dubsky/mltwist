package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
)

type Instruction struct {
	i *instruction
}

func (i Instruction) String() string { return i.i.Instr.Details.String() }

// Address returns the in-memory address of the instruction in the original
// binary.
func (i Instruction) Address() model.Address { return i.i.Instr.Address }

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
