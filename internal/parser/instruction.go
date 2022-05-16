package parser

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Instruction struct {
	Type model.Type

	Addr  model.Addr
	Bytes []byte

	Effects []expr.Effect
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
