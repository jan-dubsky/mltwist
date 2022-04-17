package repr

import (
	"decomp/pkg/expr"
	"decomp/pkg/model"
)

type Instruction struct {
	model.Instruction

	Address model.Address
	Bytes   []byte

	InputRegs  RegSet
	OutputRegs RegSet

	jumpTargets []expr.Expr
}

func NewInstruction(ins model.Instruction, addr model.Address, bytes []byte) Instruction {
	inRegs, outRegs := regs(ins.Effects)
	return Instruction{
		Instruction: ins,
		Address:     addr,
		Bytes:       bytes,

		InputRegs:  inRegs,
		OutputRegs: outRegs,

		jumpTargets: jumps(ins.Effects),
	}
}
