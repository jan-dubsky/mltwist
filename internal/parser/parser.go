package parser

import (
	"fmt"
	"mltwist/internal/memory"
	"mltwist/internal/repr"
	"mltwist/pkg/model"
)

type Program struct {
	Entrypoint   model.Addr
	Memory       *memory.Memory
	Instructions []repr.Instruction
}

func Parse(
	entrypoint model.Addr,
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
			addr += model.Addr(ins.ByteLen)
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
	addr model.Addr,
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
