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

	Effects     []expr.Effect
	JumpTargets []expr.Expr
}

func newInstruction(ins parser.Instruction) Instruction {
	return Instruction{
		Type:    ins.Type,
		Addr:    ins.Addr,
		Bytes:   ins.Bytes,
		Details: ins.Details,

		Effects:     ins.Effects,
		JumpTargets: jumps(ins),
	}
}

func convertInstructions(instrs []parser.Instruction) []Instruction {
	instructions := make([]Instruction, len(instrs))
	for i, ins := range instrs {
		instructions[i] = newInstruction(ins)
	}

	return instructions
}

func jumps(ins parser.Instruction) []expr.Expr {
	var jumpAddrs []expr.Expr
	for _, ef := range ins.Effects {
		e, ok := ef.(expr.RegStore)
		if !ok {
			continue
		}

		if e.Key() != expr.IPKey {
			continue
		}

		addrs := exprtransform.JumpAddrs(e.Value())

		// Filter those jump addresses which jump to the following
		// instruction as those are technically not jumps.
		j := 0
		for i := 0; i < len(addrs); i, j = i+1, j+1 {
			addrs[j] = addrs[i]

			c, ok := addrs[i].(expr.Const)
			if !ok {
				continue
			}

			addr, _ := expr.ConstUint[model.Addr](c)
			if addr != ins.NextAddr() {
				continue
			}

			j--
		}
		jumpAddrs = append(jumpAddrs, addrs[:j]...)
	}

	return jumpAddrs
}

func (i Instruction) Len() model.Addr { return model.Addr(len(i.Bytes)) }

// NextAddr returns address following this instruction.
func (i Instruction) NextAddr() model.Addr { return i.Addr + i.Len() }
