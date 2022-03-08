package parser

import (
	"decomp/internal/memory"
	"decomp/internal/repr"
	"decomp/pkg/model"
	"fmt"
)

type Program struct {
	Entrypoint   model.Address
	Memory       *memory.Memory
	Instructions []repr.Instruction
}

func Parse(
	entrypoint model.Address,
	m *memory.Memory,
	s Parser,
) (Program, error) {
	instrs := make([]repr.Instruction, 0)
	for _, block := range m.Blocks {
		for addr := block.Begin(); addr < block.End(); {
			b := block.Addr(addr)

			ins, err := s.Parse(b)
			if err != nil {
				return Program{}, fmt.Errorf(
					"cannot parse instruction at offset 0x%x: %w",
					addr, err)
			}

			instrs = append(instrs, instrRepr(ins, addr, b))
			addr += model.Address(ins.ByteLen)
		}
	}

	return Program{
		Entrypoint:   entrypoint,
		Memory:       m,
		Instructions: instrs,
	}, nil
}

func instrRepr(
	ins model.Instruction,
	addr model.Address,
	b []byte,
) repr.Instruction {
	return repr.Instruction{
		Instruction: ins,

		Address: addr,
		Bytes:   b[:ins.ByteLen],
	}
}
