package memory_test

import (
	"mltwist/internal/memory"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlock_Addr(t *testing.T) {
	tests := []struct {
		name   string
		block  memory.Block
		addr   model.Address
		length int
	}{
		{
			name:   "block_beginning",
			block:  memory.NewBlock(64, make([]byte, 32)),
			addr:   64,
			length: 32,
		},
		{
			name:   "block_last_byte",
			block:  memory.NewBlock(64, make([]byte, 32)),
			addr:   64 + 31,
			length: 1,
		},
		{
			name:  "end",
			block: memory.NewBlock(64, make([]byte, 32)),
			addr:  64 + 32,
		},
		{
			name:  "before_begin",
			block: memory.NewBlock(64, make([]byte, 32)),
			addr:  63,
		},
		{
			name:   "middle_of_block",
			block:  memory.NewBlock(64, make([]byte, 32)),
			addr:   64 + 12,
			length: 20,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			b := tt.block.Addr(tt.addr)
			r.Equal(tt.length, len(b))
			if tt.length == 0 {
				r.Nil(b)
			}
		})
	}
}
