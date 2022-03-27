package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func testReprIns(
	tp model.Type,
	address model.Address,
	bytes model.Address,
	jmp ...model.Address,
) repr.Instruction {
	return repr.Instruction{
		Address: address,
		Instruction: model.Instruction{
			Type:        tp,
			ByteLen:     bytes,
			JumpTargets: jmp,
		},
	}
}

func TestProgram_New(t *testing.T) {
	tests := []struct {
		name   string
		seq    []repr.Instruction
		blocks []int
	}{
		{
			name: "single_block",
			seq: []repr.Instruction{
				testReprIns(model.TypeAritm, 58, 2),
				testReprIns(model.TypeAritm, 60, 3),
				testReprIns(model.TypeAritm, 63, 4),
			},
			blocks: []int{3},
		},
		{
			name: "multiple_blocks",
			seq: []repr.Instruction{
				testReprIns(model.TypeAritm, 128, 4),
				testReprIns(model.TypeAritm, 132, 4),
				testReprIns(model.TypeJump, 136, 4, 140),

				testReprIns(model.TypeAritm, 140, 2),
				testReprIns(model.TypeAritm, 142, 2),
				testReprIns(model.TypeAritm, 144, 2),
				testReprIns(model.TypeAritm, 146, 8),

				testReprIns(model.TypeAritm, 154, 2),
				testReprIns(model.TypeCJump, 156, 4, 154),
			},
			blocks: []int{3, 4, 2},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p, err := NewProgram(tt.seq)
			r.NoError(err)

			r.Equal(len(tt.blocks), p.Len())
			r.Equal(len(tt.blocks), len(p.Blocks()))
			r.Equal(len(tt.seq), p.NumInstr())
			for i, b := range p.Blocks() {
				r.Equal(b, p.Index(i))
				r.Equal(tt.blocks[i], b.Len())
			}
		})
	}
}

func TestProgram_Move(t *testing.T) {
	const numBlocks = 10

	tests := []struct {
		from   int
		to     int
		hasErr bool
		order  []int
	}{
		{
			from:  5,
			to:    5,
			order: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			from:  0,
			to:    1,
			order: []int{1, 0, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			from:  9,
			to:    1,
			order: []int{0, 9, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			from:   10,
			to:     1,
			hasErr: true,
			order:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			from:   -1,
			to:     5,
			hasErr: true,
			order:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("move_%d", i), func(t *testing.T) {
			blocks := make([]*block, numBlocks)
			for i := range blocks {
				ins := testReprIns(
					model.TypeAritm,
					model.Address(i),
					1,
				)
				blocks[i] = newBlock(i, []repr.Instruction{ins})
			}

			r := require.New(t)
			p := &Program{blocks: blocks}
			r.Equal(numBlocks, p.Len())

			err := p.Move(tt.from, tt.to)
			if tt.hasErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}

			r.Len(tt.order, numBlocks)
			for i, n := range tt.order {
				b := p.Index(i)
				r.Equal(model.Address(n), b.Begin(),
					"Invalid block at index: %d", i)
				r.Equal(i, b.Idx())
			}
		})
	}
}
