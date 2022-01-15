package riscv

import "decomp/internal/opcode"

type instr struct {
	name   string
	opcode opcode.Opcode
}

func (i *instr) Opcode() opcode.Opcode { return i.opcode }
func (i instr) String() string         { return i.name }
