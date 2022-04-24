package parser

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Instruction struct {
	Type model.Type

	Address model.Addr
	Bytes   []byte

	Effects []expr.Effect
	Details model.PlatformDetails
}

func newInstruction(ins model.Instruction, addr model.Addr, bytes []byte) Instruction {
	return Instruction{
		Type:    ins.Type,
		Address: addr,
		Bytes:   bytes[:ins.ByteLen],

		Effects: exprtransform.EffectsApply(ins.Effects, exprtransform.ConstFold),
		Details: ins.Details,
	}
}
func (i Instruction) Len() model.Addr { return model.Addr(len(i.Bytes)) }
