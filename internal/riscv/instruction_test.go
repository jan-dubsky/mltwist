package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstructionBytes_regNum(t *testing.T) {
	tests := []struct {
		name  string
		bytes InstrBytes
		r     reg
		want  regNum
	}{
		{
			name:  "rd",
			bytes: InstrBytes{0x80, 0b11100},
			r:     rd,
			want:  0b11001,
		},
		{
			name:  "rs1",
			bytes: InstrBytes{0, 0x80, 0b11100},
			r:     rs1,
			want:  0b11001,
		},
		{
			name:  "rs2",
			bytes: InstrBytes{0, 0, 0b0011 << 4, 0b11},
			r:     rs2,
			want:  0b10011,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			num := tt.bytes.regNum(tt.r)
			require.Equal(t, tt.want, num)
		})
	}
}
