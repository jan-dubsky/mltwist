package elf

import (
	"mltwist/pkg/model"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemory_New(t *testing.T) {
	tests := []struct {
		name     string
		blocks   []Block
		hasErr   bool
		blockCnt int
	}{
		{
			name: "single_block",
			blocks: []Block{
				newBlock(50, make([]byte, 42)),
			},
			blockCnt: 1,
		},
		{
			name: "two_blocks",
			blocks: []Block{
				newBlock(50, make([]byte, 42)),
				newBlock(100, make([]byte, 30)),
			},
			blockCnt: 2,
		},
		{
			name: "two_blocks_unsorted",
			blocks: []Block{
				newBlock(100, make([]byte, 30)),
				newBlock(50, make([]byte, 42)),
			},
			blockCnt: 2,
		},
		{
			name: "overlapping_blocks",
			blocks: []Block{
				newBlock(100, make([]byte, 30)),
				newBlock(50, make([]byte, 52)),
			},
			hasErr: true,
		},
		{
			name: "blocks_touching_one_another",
			blocks: []Block{
				newBlock(102, make([]byte, 30)),
				newBlock(50, make([]byte, 52)),
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

			m, err := newMemory(tt.blocks)
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

func TestMemory_Address(t *testing.T) {
	tests := []struct {
		name   string
		blocks []Block
		addr   model.Addr
		length int
	}{
		{
			name: "full_block",
			blocks: []Block{
				newBlock(52, make([]byte, 48)),
				newBlock(120, make([]byte, 60)),
				newBlock(240, make([]byte, 20)),
			},
			addr:   52,
			length: 48,
		},
		{
			name: "out_of_blocks",
			blocks: []Block{
				newBlock(52, make([]byte, 48)),
				newBlock(120, make([]byte, 60)),
				newBlock(240, make([]byte, 20)),
			},
			addr:   110,
			length: 0,
		},
		{
			name: "middle_of_block",
			blocks: []Block{
				newBlock(52, make([]byte, 48)),
				newBlock(120, make([]byte, 60)),
				newBlock(240, make([]byte, 20)),
			},
			addr:   145,
			length: 35,
		},
		{
			name: "in_last_block",
			blocks: []Block{
				newBlock(52, make([]byte, 48)),
				newBlock(120, make([]byte, 60)),
				newBlock(240, make([]byte, 20)),
			},
			addr:   250,
			length: 10,
		},
		{
			name: "behind_the_last_block",
			blocks: []Block{
				newBlock(52, make([]byte, 48)),
				newBlock(120, make([]byte, 60)),
				newBlock(240, make([]byte, 20)),
			},
			addr:   300,
			length: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			m, err := newMemory(tt.blocks)
			r.NoError(err)
			r.NotNil(m)

			b := m.Address(tt.addr)
			r.Equal(tt.length, len(b))
			if tt.length == 0 {
				r.Nil(b)
			}
		})
	}
}
