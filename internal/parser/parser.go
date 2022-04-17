package parser

import (
	"mltwist/internal/memory"
	"mltwist/internal/repr"
	"mltwist/pkg/model"
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
	p Parser,
) (Program, error) {
	instrs := make([]repr.Instruction, 0, len(m.Blocks))
	for _, block := range m.Blocks {
		for addr := block.Begin(); addr < block.End(); {
			ins, err := parseIns(p, block, addr)
			if err != nil {
				return Program{}, fmt.Errorf(
					"cannot parse instruction at address 0x%x: %w",
					addr, err)
			}

			instrs = append(instrs, ins)
			addr += model.Address(ins.ByteLen)
		}
	}

	return Program{
		Entrypoint:   entrypoint,
		Memory:       m,
		Instructions: instrs,
	}, nil
}

func parseIns(
	p Parser,
	block memory.Block,
	addr model.Address,
) (repr.Instruction, error) {
	b := block.Addr(addr)

	ins, err := p.Parse(addr, b)
	if err != nil {
		return repr.Instruction{}, fmt.Errorf("parsing error: %w", err)
	}

	if err := ins.Validate(); err != nil {
		err = fmt.Errorf("invalid instruction model produced: %w", err)
		return repr.Instruction{}, err
	}

	return repr.NewInstruction(ins, addr, b), nil
}
