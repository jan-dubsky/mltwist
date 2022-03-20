package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReg_regNum(t *testing.T) {
	tests := []struct {
		name  string
		value uint32
		r     reg
		want  regNum
	}{
		{
			name:  "rd",
			value: valueFromBytes(0x80, 0b11100),
			r:     rd,
			want:  0b11001,
		},
		{
			name:  "rs1",
			value: valueFromBytes(0, 0x80, 0b11100),
			r:     rs1,
			want:  0b11001,
		},
		{
			name:  "rs2",
			value: valueFromBytes(0, 0, 0b0011<<4, 0b11),
			r:     rs2,
			want:  0b10011,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			num := tt.r.regNum(tt.value)
			require.Equal(t, tt.want, num)
		})
	}
}
