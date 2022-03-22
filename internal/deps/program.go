package deps

import (
	"decomp/internal/deps/basicblock"
	"decomp/internal/repr"
	"fmt"
)

type Program struct {
	blocks   []*Block
	instrCnt int
}

func NewProgram(seq []repr.Instruction) (*Program, error) {
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

	return &Program{
		blocks:   blocks,
		instrCnt: instrCnt,
	}, nil
}

func (m *Program) Blocks() []*Block   { return m.blocks }
func (m *Program) Len() int           { return len(m.blocks) }
func (m *Program) NumInstr() int      { return m.instrCnt }
func (m *Program) Index(i int) *Block { return m.blocks[i] }
