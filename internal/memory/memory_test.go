package memory_test

import (
	"decomp/internal/addr"
	"decomp/internal/memory"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemory_New(t *testing.T) {
	tests := []struct {
		name     string
		blocks   []memory.Block
		hasErr   bool
		blockCnt int
	}{
		{
			name:     "single_block",
			blocks:   []memory.Block{memory.NewBlock(50, make([]byte, 42))},
			blockCnt: 1,
		},
		{
			name: "two_blocks",
			blocks: []memory.Block{
				memory.NewBlock(50, make([]byte, 42)),
				memory.NewBlock(100, make([]byte, 30)),
			},
			blockCnt: 2,
		},
		{
			name: "two_blocks_unsorted",
			blocks: []memory.Block{
				memory.NewBlock(100, make([]byte, 30)),
				memory.NewBlock(50, make([]byte, 42)),
			},
			blockCnt: 2,
		},
		{
			name: "overlapping_blocks",
			blocks: []memory.Block{
				memory.NewBlock(100, make([]byte, 30)),
				memory.NewBlock(50, make([]byte, 52)),
			},
			hasErr: true,
		},
		{
			name: "blocks_touching_one_another",
			blocks: []memory.Block{
				memory.NewBlock(102, make([]byte, 30)),
				memory.NewBlock(50, make([]byte, 52)),
			},
			blockCnt: 2,
		},
		{
			name:     "no_block",
			blockCnt: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			m, err := memory.New(tt.blocks...)
			if tt.hasErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(m)
			r.Equal(tt.blockCnt, len(m.Blocks))

			r.True(sort.SliceIsSorted(m.Blocks, func(i, j int) bool {
				return m.Blocks[i].Begin() < m.Blocks[j].Begin()
			}))
		})
	}
}

func TestMemory_Addr(t *testing.T) {
	tests := []struct {
		name   string
		blocks []memory.Block
		addr   addr.Address
		length int
	}{
		{
			name: "full_block",
			blocks: []memory.Block{
				memory.NewBlock(52, make([]byte, 48)),
				memory.NewBlock(120, make([]byte, 60)),
				memory.NewBlock(240, make([]byte, 20)),
			},
			addr:   52,
			length: 48,
		},
		{
			name: "out_of_blocks",
			blocks: []memory.Block{
				memory.NewBlock(52, make([]byte, 48)),
				memory.NewBlock(120, make([]byte, 60)),
				memory.NewBlock(240, make([]byte, 20)),
			},
			addr:   110,
			length: 0,
		},
		{
			name: "middle_of_block",
			blocks: []memory.Block{
				memory.NewBlock(52, make([]byte, 48)),
				memory.NewBlock(120, make([]byte, 60)),
				memory.NewBlock(240, make([]byte, 20)),
			},
			addr:   145,
			length: 35,
		},
		{
			name: "in_last_block",
			blocks: []memory.Block{
				memory.NewBlock(52, make([]byte, 48)),
				memory.NewBlock(120, make([]byte, 60)),
				memory.NewBlock(240, make([]byte, 20)),
			},
			addr:   250,
			length: 10,
		},
		{
			name: "behind_the_last_block",
			blocks: []memory.Block{
				memory.NewBlock(52, make([]byte, 48)),
				memory.NewBlock(120, make([]byte, 60)),
				memory.NewBlock(240, make([]byte, 20)),
			},
			addr:   300,
			length: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			m, err := memory.New(tt.blocks...)
			r.NoError(err)
			r.NotNil(m)

			b := m.Addr(tt.addr)
			r.Equal(tt.length, len(b))
			if tt.length == 0 {
				r.Nil(b)
			}
		})
	}
}
