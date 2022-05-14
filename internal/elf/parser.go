package elf

import (
	"debug/elf"
	"fmt"
	"io"
	"mltwist/pkg/model"
)

type Parser struct {
	f *elf.File
}

func NewParser(filename string) (*Parser, error) {
	f, err := elf.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", filename, err)
	}

	if f.Type&elf.ET_EXEC == 0 {
		return nil, fmt.Errorf("file %q is not an executable ELF file", filename)
	}

	return &Parser{f: f}, nil
}

func (p *Parser) Entrypoint() model.Addr { return model.Addr(p.f.Entry) }

func (p *Parser) MachineCode() (*Memory, error) {
	var blocks []Block
	for _, s := range p.f.Sections {
		if skipMachineCodeSection(s) {
			continue
		}
		data, err := s.Data()
		if err != nil {
			return nil, fmt.Errorf("cannot read section %q: %w", s.Name, err)
		}

		if uint64(len(data)) != s.Size {
			return nil, fmt.Errorf(
				"size of ELF section and read bytes differ: %d != %d",
				s.Size, len(data))
		}

		b := newBlock(model.Addr(s.Addr), data)
		blocks = append(blocks, b)
	}

	return nonEmptyMemory(blocks)
}

func (p *Parser) Memory() (*Memory, error) {
	var blocks []Block
	for _, p := range p.f.Progs {
		if p.Type != elf.PT_LOAD {
			continue
		}

		if p.Memsz < p.Filesz {
			return nil, fmt.Errorf(
				"program section in memory less then in file: %d < %d",
				p.Memsz, p.Filesz)
		}

		data, err := io.ReadAll(p.Open())
		if err != nil {
			return nil, fmt.Errorf("cannot read program section: %w", err)
		}

		if missing := p.Memsz - uint64(len(data)); missing > 0 {
			data = append(data, make([]byte, missing)...)
		}

		b := newBlock(model.Addr(p.Vaddr), data)
		blocks = append(blocks, b)
	}

	return nonEmptyMemory(blocks)
}

func nonEmptyMemory(blocks []Block) (*Memory, error) {
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no non-empty memory blocks found")
	}

	mem, err := newMemory(blocks)
	if err != nil {
		return nil, fmt.Errorf("memory creation failed: %w", err)
	}

	return mem, nil
}

func skipMachineCodeSection(s *elf.Section) bool {
	if s.Type != elf.SHT_PROGBITS || s.Size == 0 {
		return true
	}

	// man elf states: "If this section appears in the memory image of a
	// process, this member holds the address at which the section's first
	// byte should reside.  Otherwise, the member contains zero".
	// Consequentlyzero address is not represented in program memory address
	// space.
	if s.Addr == 0 {
		return true
	}

	if s.Flags&elf.SHF_EXECINSTR == 0 {
		return true
	}

	return false
}

func (p *Parser) Close() error {
	err := p.f.Close()
	if err != nil {
		return fmt.Errorf("elf file close failed: %w", err)
	}

	p.f = nil
	return nil
}
