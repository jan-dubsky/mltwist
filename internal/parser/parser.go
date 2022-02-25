package parser

import (
	"decomp/internal/addr"
	"decomp/internal/instruction"
	"decomp/internal/memory"
	"fmt"
)

type Program struct {
	Entrypoint   addr.Address
	Memory       *memory.Memory
	Instructions []instruction.Instruction
}

func Parse(
	entrypoint addr.Address,
	m *memory.Memory,
	s Strategy,
) (Program, error) {
	instrs := make([]instruction.Instruction, 0)
	for _, block := range m.Blocks {
		for addr := block.Begin(); addr < block.End(); {
			b := block.Addr(addr)

			instr, err := s.Parse(b)
			if err != nil {
				return Program{}, fmt.Errorf(
					"cannot parse instruction at offset 0x%x: %w",
					addr, err)
			}

			instr.Bytes = b[:instr.ByteLen]
			instr.Address = addr

			instrs = append(instrs, instr)
			addr += instr.ByteLen
		}
	}

	return Program{
		Entrypoint:   entrypoint,
		Memory:       m,
		Instructions: instrs,
	}, nil
}
