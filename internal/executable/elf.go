package executable

import (
	"debug/elf"
	"fmt"
	"io"
	"mltwist/internal/memory"
	"mltwist/pkg/model"
)

type Executable struct {
	Entrypoint model.Addr
	Memory     *memory.Memory
}

func Parse(filename string) (Executable, error) {
	f, err := elf.Open(filename)
	if err != nil {
		return Executable{}, fmt.Errorf("cannot open file %q: %w", filename, err)
	}
	defer f.Close()

	blocks := make([]memory.Block, 0, 1)
	for _, s := range f.Sections {
		if s.Size == 0 {
			continue
		}

		if s.Flags&elf.SHF_EXECINSTR == 0 {
			continue
		}

		if s.Flags&elf.SHF_COMPRESSED != 0 {
			err = fmt.Errorf("section %s is compressed", s.Name)
			return Executable{}, err
		}

		bytes, err := io.ReadAll(s.Open())
		if err != nil {
			err = fmt.Errorf("cannot read bytes of section %q", s.Name)
			return Executable{}, err
		}

		if uint64(len(bytes)) != s.Size {
			return Executable{}, fmt.Errorf(
				"size of ELF section and read bytes differ: %d != %d",
				s.Size, len(bytes))
		}

		blocks = append(blocks, memory.NewBlock(model.Addr(s.Addr), bytes))
	}

	if len(blocks) == 0 {
		return Executable{}, fmt.Errorf("no executable blocks found")
	}

	mem, err := memory.New(blocks...)
	if err != nil {
		return Executable{}, fmt.Errorf("memory creation failed: %w", err)
	}

	return Executable{
		Entrypoint: model.Addr(f.Entry),
		Memory:     mem,
	}, nil
}
