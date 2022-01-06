package executable

import (
	"fmt"
	"sort"
)

type codeBlock struct {
	begin uint64
	bytes []byte
}

func (b *codeBlock) end() uint64 { return b.begin + uint64(len(b.bytes)) }

func (b *codeBlock) slice(pos uint64, length uint64) []byte {
	offset := pos - b.begin
	end := offset + length
	if end > uint64(len(b.bytes)) {
		return nil
	}

	return b.bytes[offset:end]
}

type Executable struct {
	Entrypoint uint64
	blocks     []codeBlock
}

func (c *Executable) Bytes(pos uint64, length uint64) []byte {
	idx := sort.Search(len(c.blocks), func(i int) bool {
		return c.blocks[i].end() > pos
	})

	if idx == len(c.blocks) {
		return nil
	}

	return c.blocks[idx].slice(pos, length)
}

func newExecutable(entrypoint uint64, blocks []codeBlock) (Executable, error) {
	if len(blocks) == 0 {
		return Executable{}, fmt.Errorf("no executable blocks found")
	}

	sort.Slice(blocks, func(i, j int) bool { return blocks[i].begin < blocks[j].begin })
	for i, b := range blocks {
		if i == 0 {
			continue
		} else if prev := blocks[i-1]; prev.end() <= b.begin {
			continue
		}

		return Executable{}, fmt.Errorf("code sections overlap one another")
	}

	return Executable{
		Entrypoint: entrypoint,
		blocks:     blocks,
	}, nil
}
