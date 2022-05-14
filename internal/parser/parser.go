package parser

import (
	"fmt"
	"mltwist/internal/elf"
	"mltwist/pkg/model"
)

// Parse parses all instructions in a memory. This function fails if any of
// parsings fails.
func Parse(
	m *elf.Memory,
	p Parser,
) ([]Instruction, error) {
	instrs := make([]Instruction, 0, len(m.Blocks))
	for _, block := range m.Blocks {
		for addr := block.Begin(); addr < block.End(); {
			ins, err := parseIns(p, block, addr)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot parse instruction at address 0x%x: %w",
					addr, err)
			}

			instrs = append(instrs, ins)
			addr += ins.Len()
		}
	}

	return instrs, nil
}

func parseIns(
	p Parser,
	block elf.Block,
	addr model.Addr,
) (Instruction, error) {
	b := block.Address(addr)

	ins, err := p.Parse(addr, b)
	if err != nil {
		return Instruction{}, fmt.Errorf("parsing error: %w", err)
	}

	if err := ins.Validate(); err != nil {
		err = fmt.Errorf("invalid instruction model produced: %w", err)
		return Instruction{}, err
	}

	return newInstruction(ins, addr, b), nil
}
