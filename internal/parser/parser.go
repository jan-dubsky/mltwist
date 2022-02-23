package parser

import (
	"decomp/internal/instruction"
	"decomp/internal/memory"
	"fmt"
)

type Instructions struct {
	Entrypoint   memory.Address
	Memory       *memory.Memory
	Instructions []instruction.Instruction
}

func Parse(
	entrypoint memory.Address,
	m *memory.Memory,
	s Strategy,
) (Instructions, error) {
	instrs := make([]instruction.Instruction, 0)
	for _, block := range m.Blocks {
		for addr := block.Begin; addr < block.End(); {
			instr, err := s.Parse(block.Addr(addr))
			if err != nil {
				return Instructions{}, fmt.Errorf(
					"cannot parse instruction at offset 0x%x: %w",
					addr, err)
			}

			instrs = append(instrs, instr)
			fmt.Printf("%d (0x%x): %s\n", len(instrs), addr, instr.Details.String())
			addr += instr.ByteLen
		}
	}

	return Instructions{
		Entrypoint:   entrypoint,
		Memory:       m,
		Instructions: instrs,
	}, nil
}
