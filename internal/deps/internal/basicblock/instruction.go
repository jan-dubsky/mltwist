package basicblock

import (
	"mltwist/internal/exprtransform"
	"mltwist/internal/parser"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Instruction struct {
	Type    model.Type
	Addr    model.Addr
	Bytes   []byte
	Details model.PlatformDetails

	Effs        []expr.Effect
	JumpTargets []expr.Expr
}

func newInstruction(ins parser.Instruction) Instruction {
	return Instruction{
		Type:    ins.Type,
		Addr:    ins.Address,
		Bytes:   ins.Bytes,
		Details: ins.Details,

		Effs:        ins.Effects,
		JumpTargets: jumps(ins.Effects),
	}
}

func convertInstructions(instrs []parser.Instruction) []Instruction {
	instructions := make([]Instruction, len(instrs))
	for i, ins := range instrs {
		instructions[i] = newInstruction(ins)
	}

	return instructions
}

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

func (i Instruction) Len() model.Addr { return model.Addr(len(i.Bytes)) }

// NextAddr returns address following this instruction.
func (i Instruction) NextAddr() model.Addr { return i.Addr + i.Len() }

func (i Instruction) Effects() []expr.Effect { return i.Effs }
