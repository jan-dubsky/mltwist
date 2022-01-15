package parser

import "fmt"

type parsing struct {
	m      Memory
	s      Strategy
	window uint64
}

type Instructions struct{}

func Parse(entrypoint uint64, m Memory, s Strategy) (Instructions, error) {
	p := parsing{
		m:      m,
		s:      s,
		window: s.Window(),
	}

	err := p.parse(entrypoint)
	if err != nil {
		return Instructions{}, err
	}

	return Instructions{}, nil
}

func (p *parsing) parse(offset uint64) error {
	bytes := p.m.Bytes(offset, p.window)
	if bytes == nil {
		return fmt.Errorf("cannot read %d bytes at offset 0x%x", p.window, offset)
	}

	instr, err := p.s.Parse(bytes)
	if err != nil {
		return fmt.Errorf("cannot parse instruction ot offset 0x%x: %w",
			offset, err)
	}

	fmt.Printf("instruction: %s\n", instr.Details.String())

	// Jump instructions unconditionally jump elsewhere, which makes
	// instruction following jump instruction unreachable (unless addressed
	// by another jump/CJump/call/... instruction).
	if !instr.Type.Jump() {
		err = p.parse(offset + instr.ByteLen)
		if err != nil {
			return err
		}
	}

	for _, t := range instr.JumpTargets {
		err = p.parse(uint64(t))
		if err != nil {
			return err
		}
	}

	return nil
}
