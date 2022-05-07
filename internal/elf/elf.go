package elf

import (
	"debug/elf"
	"fmt"
	"mltwist/pkg/model"
)

type Executable struct {
	Entrypoint model.Addr
	Memory     *Memory
}

func Parse(filename string) (Executable, error) {
	f, err := elf.Open(filename)
	if err != nil {
		return Executable{}, fmt.Errorf("cannot open file %q: %w", filename, err)
	}
	defer f.Close()

	if f.Type&elf.ET_EXEC == 0 {
		return Executable{}, fmt.Errorf("file %q is not executable", filename)
	}

	// ELF file starts by zero section which we will skip ,that's why we
	// preallocate just len(f.Sections) - 1 blocks.
	blocks := make([]Block, 0, len(f.Sections)-1)
	for _, s := range f.Sections {
		if s.Size == 0 || s.Addr == 0 {
			continue
		}

		data, err := s.Data()
		if err != nil {
			err = fmt.Errorf("cannot read section %q: %w", s.Name, err)
			return Executable{}, err
		}

		if uint64(len(data)) != s.Size {
			return Executable{}, fmt.Errorf(
				"size of ELF section and read bytes differ: %d != %d",
				s.Size, len(data))
		}

		b := newBlock(model.Addr(s.Addr), data, s.Flags&elf.SHF_EXECINSTR != 0)
		blocks = append(blocks, b)
	}

	if len(blocks) == 0 {
		return Executable{}, fmt.Errorf("no non-empty blocks found")
	}

	mem, err := newMemory(blocks)
	if err != nil {
		return Executable{}, fmt.Errorf("memory creation failed: %w", err)
	}

	return Executable{
		Entrypoint: model.Addr(f.Entry),
		Memory:     mem,
	}, nil
}
