package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
)

type Instruction struct {
	Address model.Address
	Instr   repr.Instruction

	deps []*Instruction
}

func newInstruction(ins repr.Instruction) *Instruction {
	return &Instruction{
		Address: ins.Address,
		Instr:   ins,
	}
}

func (i *Instruction) LowerBound() model.Address {
	var max *Instruction
	for _, d := range i.deps {
		if max == nil || d.Address > max.Address {
			max = d
		}
	}

	if max == nil {
		return model.MinAddress
	}

	return max.Address + max.Instr.ByteLen
}
