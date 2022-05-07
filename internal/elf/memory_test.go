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
				newBlock(50, make([]byte, 42), true),
			},
			blockCnt: 1,
		},
		{
			name: "two_blocks",
			blocks: []Block{
				newBlock(50, make([]byte, 42), true),
				newBlock(100, make([]byte, 30), true),
			},
			blockCnt: 2,
		},
		{
			name: "two_blocks_unsorted",
			blocks: []Block{
				newBlock(100, make([]byte, 30), true),
				newBlock(50, make([]byte, 42), true),
			},
			blockCnt: 2,
		},
		{
			name: "overlapping_blocks",
			blocks: []Block{
				newBlock(100, make([]byte, 30), true),
				newBlock(50, make([]byte, 52), true),
			},
			hasErr: true,
		},
		{
			name: "blocks_touching_one_another",
			blocks: []Block{
				newBlock(102, make([]byte, 30), true),
				newBlock(50, make([]byte, 52), true),
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
				newBlock(52, make([]byte, 48), true),
				newBlock(120, make([]byte, 60), true),
				newBlock(240, make([]byte, 20), true),
			},
			addr:   52,
			length: 48,
		},
		{
			name: "out_of_blocks",
			blocks: []Block{
				newBlock(52, make([]byte, 48), true),
				newBlock(120, make([]byte, 60), true),
				newBlock(240, make([]byte, 20), true),
			},
			addr:   110,
			length: 0,
		},
		{
			name: "middle_of_block",
			blocks: []Block{
				newBlock(52, make([]byte, 48), true),
				newBlock(120, make([]byte, 60), true),
				newBlock(240, make([]byte, 20), true),
			},
			addr:   145,
			length: 35,
		},
		{
			name: "in_last_block",
			blocks: []Block{
				newBlock(52, make([]byte, 48), true),
				newBlock(120, make([]byte, 60), true),
				newBlock(240, make([]byte, 20), true),
			},
			addr:   250,
			length: 10,
		},
		{
			name: "behind_the_last_block",
			blocks: []Block{
				newBlock(52, make([]byte, 48), true),
				newBlock(120, make([]byte, 60), true),
				newBlock(240, make([]byte, 20), true),
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
