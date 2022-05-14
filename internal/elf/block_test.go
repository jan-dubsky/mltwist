package elf

import (
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlock_Address(t *testing.T) {
	tests := []struct {
		name   string
		block  Block
		addr   model.Addr
		length int
	}{
		{
			name:   "block_beginning",
			block:  newBlock(64, make([]byte, 32)),
			addr:   64,
			length: 32,
		},
		{
			name:   "block_last_byte",
			block:  newBlock(64, make([]byte, 32)),
			addr:   64 + 31,
			length: 1,
		},
		{
			name:  "end",
			block: newBlock(64, make([]byte, 32)),
			addr:  64 + 32,
		},
		{
			name:  "before_begin",
			block: newBlock(64, make([]byte, 32)),
			addr:  63,
		},
		{
			name:   "middle_of_block",
			block:  newBlock(64, make([]byte, 32)),
			addr:   64 + 12,
			length: 20,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			b := tt.block.Address(tt.addr)
			r.Equal(tt.length, len(b))
			if tt.length == 0 {
				r.Nil(b)
			}
		})
	}
}
