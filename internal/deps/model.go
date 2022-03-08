package deps

import (
	"decomp/internal/deps/basicblock"
	"decomp/internal/repr"
	"fmt"
)

type Model struct {
	blocks []*Block
}

func NewModel(seq []repr.Instruction) (*Model, error) {
	seqs, err := basicblock.Parse(seq)
	if err != nil {
		return nil, fmt.Errorf("cannot find basic blocks: %w", err)
	}

	blocks := make([]*Block, len(seqs))
	for i, seq := range seqs {
		blocks[i] = newBlock(seq)
	}

	return &Model{
		blocks: blocks,
	}, nil
}

func (m *Model) Blocks() []*Block { return m.blocks }
