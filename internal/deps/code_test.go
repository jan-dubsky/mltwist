package deps

import (
	"fmt"
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/parser"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func testInputInsJump(
	addr model.Addr,
	bytes model.Addr,
	jumpAddrs ...model.Addr,
) parser.Instruction {
	jumps := make([]expr.Effect, len(jumpAddrs))
	for i, j := range jumpAddrs {
		jumps[i] = expr.NewRegStore(model.AddrExpr(j), expr.IPKey, expr.Width32)
	}

	return parser.Instruction{
		Addr:    addr,
		Bytes:   make([]byte, bytes),
		Effects: jumps,
	}
}

func TestCode_New(t *testing.T) {
	tests := []struct {
		name   string
		seq    []parser.Instruction
		blocks []int
	}{
		{
			name: "single_block",
			seq: []parser.Instruction{
				testInputInsJump(58, 2),
				testInputInsJump(60, 3),
				testInputInsJump(63, 4),
			},
			blocks: []int{3},
		},
		{
			name: "multiple_blocks",
			seq: []parser.Instruction{
				testInputInsJump(128, 4),
				testInputInsJump(132, 4),
				testInputInsJump(136, 4, 128),

				testInputInsJump(140, 2),
				testInputInsJump(142, 2),
				testInputInsJump(144, 2),
				testInputInsJump(146, 8),

				testInputInsJump(154, 2),
				testInputInsJump(156, 4, 154),
			},
			blocks: []int{3, 4, 2},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			p, err := NewCode(tt.seq[0].Addr, tt.seq)
			r.NoError(err)

			r.Equal(len(tt.blocks), p.Len())
			r.Equal(len(tt.blocks), len(p.Blocks()))
			r.Equal(len(tt.seq), p.NumInstr())
			for i, b := range p.Blocks() {
				r.Equal(b, p.Index(i))
				r.Equal(tt.blocks[i], b.Num())
			}
		})
	}
}

func TestCode_Move(t *testing.T) {
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
				ins := basicblock.Instruction{Addr: model.Addr(i)}
				blocks[i] = newBlock(i, []basicblock.Instruction{ins})
			}

			r := require.New(t)
			p := &Code{blocks: blocks}
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
				r.Equal(model.Addr(n), b.Begin(),
					"Invalid block at index: %d", i)
				r.Equal(i, b.Idx())
			}
		})
	}
}
