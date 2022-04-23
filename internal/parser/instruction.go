package parser

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Instruction struct {
	model.Instruction

	Address model.Addr
	Bytes   []byte

	JumpTargets []expr.Expr
}

func newInstruction(ins model.Instruction, addr model.Addr, bytes []byte) Instruction {
	return Instruction{
		Instruction: ins,
		Address:     addr,
		Bytes:       bytes,

		JumpTargets: jumps(ins.Effects),
	}
}

// NextAddr returns address following this instruction.
func (i Instruction) NextAddr() model.Addr { return i.Address + i.ByteLen }

func jumps(effects []expr.Effect) []expr.Expr {
	var jumpAddrs []expr.Expr
	for _, ef := range effects {
		e, ok := ef.(expr.RegStore)
		if !ok {
			continue
		}

		if e.Key() != expr.IPKey {
			continue
		}

		addrs := exprtransform.JumpAddrs(e.Value())
		jumpAddrs = append(jumpAddrs, addrs...)
	}

	return jumpAddrs
}
