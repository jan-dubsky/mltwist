package parser

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// Instruction represents single parsed machine code instruction.
type Instruction struct {
	Type model.Type

	// Addr is memory address of the instruction in program virtual memory.
	Addr model.Addr
	// Bytes is slice of raw bytes representing the instruction in program
	// memory.
	Bytes []byte

	// Effects is a constant-folded list of expression side-effects
	// representing the instruction functionality.
	Effects []expr.Effect
	// Details provide platform-dependent functionality of the instruction.
	Details model.PlatformDetails
}

func newInstruction(ins model.Instruction, addr model.Addr, bytes []byte) Instruction {
	return Instruction{
		Type:  ins.Type,
		Addr:  addr,
		Bytes: bytes[:ins.ByteLen],

		Effects: exprtransform.EffectsApply(ins.Effects, exprtransform.ConstFold),
		Details: ins.Details,
	}
}

func (i Instruction) Begin() model.Addr { return i.Addr }
func (i Instruction) Len() model.Addr   { return model.Addr(len(i.Bytes)) }
func (i Instruction) End() model.Addr   { return i.Addr + i.Len() }
