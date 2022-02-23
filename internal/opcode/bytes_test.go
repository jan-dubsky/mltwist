package opcode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyMask(t *testing.T) {
	tests := []struct {
		bytes    []byte
		mask     []byte
		expected []byte
	}{
		{
			bytes:    []byte{0xf0, 0x0f, 0x33, 0xbb},
			mask:     []byte{0x3f, 0xfb, 0x11, 0xf8},
			expected: []byte{0x30, 0x0b, 0x11, 0xb8},
		},
		{
			bytes:    []byte{0xf0, 0x0f, 0x33, 0xbb},
			mask:     []byte{0x3f, 0xfb},
			expected: []byte{0x30, 0x0b},
		},
	}

	r := require.New(t)
	for _, tt := range tests {
		actual := applyMask(tt.bytes, tt.mask)
		r.Equal(tt.expected, actual)
	}
}
