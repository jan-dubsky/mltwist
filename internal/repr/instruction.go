package repr

import (
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Instruction struct {
	model.Instruction

	Address model.Addr
	Bytes   []byte

	InputRegistry  RegSet
	OutputRegistry RegSet

	JumpTargets []expr.Expr
}

func NewInstruction(ins model.Instruction, addr model.Addr, bytes []byte) Instruction {
	inRegs, outRegs := regs(ins.Effects)
	return Instruction{
		Instruction: ins,
		Address:     addr,
		Bytes:       bytes,

		InputRegistry:  inRegs,
		OutputRegistry: outRegs,

		JumpTargets: jumps(ins.Effects),
	}
}

// NextAddr returns address following this instruction.
func (i Instruction) NextAddr() model.Addr { return i.Address + i.ByteLen }
