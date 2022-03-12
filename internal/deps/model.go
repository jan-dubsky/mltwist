package deps

import (
	"decomp/internal/deps/basicblock"
	"decomp/internal/repr"
	"fmt"
)

type Model struct {
	blocks   []*Block
	instrCnt int
}

func NewModel(seq []repr.Instruction) (*Model, error) {
	seqs, err := basicblock.Parse(seq)
	if err != nil {
		return nil, fmt.Errorf("cannot find basic blocks: %w", err)
	}

	blocks := make([]*Block, len(seqs))
	for i, seq := range seqs {
		blocks[i] = newBlock(i, seq)
	}

	var instrCnt int
	for _, b := range blocks {
		instrCnt += b.Len()
	}

	return &Model{
		blocks:   blocks,
		instrCnt: instrCnt,
	}, nil
}

func (m *Model) Blocks() []*Block   { return m.blocks }
func (m *Model) Len() int           { return len(m.blocks) }
func (m *Model) NumInstr() int      { return m.instrCnt }
func (m *Model) Index(i int) *Block { return m.blocks[i] }
