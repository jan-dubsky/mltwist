package deps

import "decomp/pkg/model"

type Instruction struct {
	i *instruction
}

func (i Instruction) String() string { return i.i.Instr.Details.String() }

// Address returns the in-memory address of the instruction in the original
// binary.
func (i Instruction) Address() model.Address { return i.i.Instr.Address }

func (i Instruction) Idx() int { return i.i.Idx() }
