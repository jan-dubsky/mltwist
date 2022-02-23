package parser

import "decomp/internal/instruction"

type Strategy interface {
	Parse(bytes []byte) (instruction.Instruction, error)
}
