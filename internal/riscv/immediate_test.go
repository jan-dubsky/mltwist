package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBitRange(t *testing.T) {
	r := require.New(t)

	b := InstrBytes{0, 0, 0xff, 0xff}
	r.Equal((uint32(1)<<12)-1, parseBitRange(b, 20, 32))

	b = InstrBytes{0b1010 << 4, 0b111000}
	r.Equal(uint32(0b10001010), parseBitRange(b, 4, 12))
}

func TestSignExtend(t *testing.T) {
	r := require.New(t)

	r.Equal(int32(5), signExtend(5, 11))
	r.Equal(int32(2047), signExtend(2047, 11))
	r.Equal(int32(-2048), signExtend(2048, 11))
	r.Equal(int32(-1), signExtend((1<<12)-1, 11))
	r.Equal(int32(-2050), signExtend(0b11011111111110, 13))
	r.Equal(int32(-1), signExtend(1, 0))
}

func TestImmediate_parseValue(t *testing.T) {
	tests := []struct {
		name          string
		immType       immType
		b             InstrBytes
		expected      int32
		expectedFalse bool
	}{
		{
			name:          "R-type",
			immType:       immTypeR,
			b:             InstrBytes{0xff, 0xff, 0xff, 0xff},
			expected:      0,
			expectedFalse: true,
		},
		{
			name:     "I-type_positive",
			immType:  immTypeI,
			b:        InstrBytes{0, 0, 0xfe, 0b00101010},
			expected: 0x2af,
		},
		{
			name:     "I-type_negative",
			immType:  immTypeI,
			b:        InstrBytes{0, 0, 0xfe, 0b10101010},
			expected: -1 - 0x550,
		},
		{
			name:     "S-type_positive",
			immType:  immTypeS,
			b:        InstrBytes{0x80, 0b1100, 0, 0b00101010},
			expected: 0x2b9,
		},
		{
			name:     "S-type_negative",
			immType:  immTypeS,
			b:        InstrBytes{0x80, 0b1100, 0, 0b10101010},
			expected: -1 - 0x546,
		},
		{
			name:     "B-type_positive",
			immType:  immTypeB,
			b:        InstrBytes{0x80, 0b0011, 0, 0b00101010},
			expected: 0x553 << 1,
		},
		{
			name:    "B-type_negative",
			immType: immTypeB,
			b:       InstrBytes{0x80, 0b0011, 0, 0b10101010},
			// As bit [0] is set to zero by hardware, it behaves
			// like if we subtracted one additional one.
			expected: -1 - (0x2ac << 1) - 1,
		},
		{
			name:     "U-type_positive",
			immType:  immTypeU,
			b:        InstrBytes{0, 0b1100 << 4, 0b10101010, 0b00110011},
			expected: 0x33aac << 12,
		},
		{
			name:    "U-type_negative",
			immType: immTypeU,
			b:       InstrBytes{0, 0b1100 << 4, 0b10101010, 0b11000011},
			// Bits [0..11] are set to 0 by hardware, so we have to
			// subtract their value here.
			expected: -1 - 0x3c553<<12 - (1<<12 - 1),
		},
		{
			name:     "J-type_positive",
			immType:  immTypeJ,
			b:        InstrBytes{0, 0b11100000, 0b00010111, 0b00110011},
			expected: 0x7eb30,
		},
		{
			name:     "J-type_negative",
			immType:  immTypeJ,
			b:        InstrBytes{0, 0b11100000, 0b00010111, 0b10110011},
			expected: -1 - 0x814cf,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			val, ok := tt.immType.parseValue(tt.b)
			if tt.expectedFalse {
				r.False(ok)
			}

			t.Logf("expected: %x\tactual: %x\n", tt.expected, val)
			r.Equal(tt.expected, val)
		})
	}
}
