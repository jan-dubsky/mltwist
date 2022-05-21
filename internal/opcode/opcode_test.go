package opcode_test

import (
	"mltwist/internal/opcode"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpcode_Validate(t *testing.T) {
	tests := []struct {
		opcode opcode.Opcode
		valid  bool
	}{{
		opcode: opcode.Opcode{
			Bytes: make([]byte, 4),
			Mask:  []byte{0, 0, 0, 1},
		},
		valid: true,
	}, {
		opcode: opcode.Opcode{
			Bytes: make([]byte, 4),
			Mask:  []byte{0, 0, 1},
		},
		valid: false,
	}, {
		opcode: opcode.Opcode{
			Bytes: make([]byte, 4),
			Mask:  []byte{0, 0, 0, 0},
		},
		valid: false,
	}, {
		opcode: opcode.Opcode{},
		valid:  false,
	}}

	r := require.New(t)
	for _, tt := range tests {
		err := tt.opcode.Validate()
		if tt.valid {
			r.NoError(err)
		} else {
			r.Error(err)
		}
	}
}
