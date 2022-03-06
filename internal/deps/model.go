package deps

import "decomp/internal/repr"

type Model struct {
	blocks []*Block
}

func NewModel(seqs [][]repr.Instruction) *Model {
	blocks := make([]*Block, len(seqs))
	for i, seq := range seqs {
		blocks[i] = newBlock(seq)
	}

	return &Model{
		blocks: blocks,
	}
}

func (m *Model) Blocks() []*Block { return m.blocks }
