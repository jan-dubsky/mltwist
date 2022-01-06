package executable

import (
	"debug/elf"
	"fmt"
	"io"
)

func Parse(filename string) (Executable, error) {
	f, err := elf.Open(filename)
	if err != nil {
		return Executable{}, fmt.Errorf("cannot open file %s: %w", filename, err)
	}
	defer f.Close()

	blocks := make([]codeBlock, 0, 1)
	for _, s := range f.Sections {
		if s.Size == 0 {
			continue
		}

		if s.Flags&elf.SHF_EXECINSTR == 0 {
			continue
		}
		if s.Flags&elf.SHF_COMPRESSED != 0 {
			return Executable{}, fmt.Errorf("section %s is compressed", s.Name)
		}

		bytes, err := io.ReadAll(s.Open())
		if err != nil {
			err = fmt.Errorf("cannot read bytes of section %s", s.Name)
			return Executable{}, err
		}

		if uint64(len(bytes)) != s.Size {
			return Executable{}, fmt.Errorf(
				"length of section and read bytes differ: %d != %d",
				s.Size, len(bytes))
		}

		b := codeBlock{
			begin: s.Addr,
			bytes: bytes,
		}
		blocks = append(blocks, b)
	}

	return newExecutable(f.Entry, blocks)
}
