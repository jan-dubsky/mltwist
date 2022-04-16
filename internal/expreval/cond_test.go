package expreval

import (
	"decomp/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEq(t *testing.T) {
	tests := []struct {
		name     string
		v1       value
		v2       value
		w        expr.Width
		expected bool
	}{
		{
			name:     "equal",
			v1:       value{5, 7, 33, 249},
			v2:       value{5, 7, 33, 249},
			w:        expr.Width32,
			expected: true,
		},
		{
			name:     "not_equal",
			v1:       valUint(34433556, expr.Width32),
			v2:       valUint(438483843, expr.Width32),
			w:        expr.Width32,
			expected: false,
		},
		{
			name:     "cut_extend",
			v1:       value{45, 135, 0, 0, 34, 67, 87},
			v2:       value{45, 135},
			w:        expr.Width32,
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := eq(tt.v1, tt.v2, tt.w)
			require.Equal(t, tt.expected, res)
		})
	}
}

func TestLtuAndLeu(t *testing.T) {
	tests := []struct {
		name string
		v1   value
		v2   value
		w    expr.Width
		ltu  bool
		leu  bool
	}{
		{
			name: "equal",
			v1:   value{5, 7, 33, 249},
			v2:   value{5, 7, 33, 249},
			w:    expr.Width32,
			ltu:  false,
			leu:  true,
		},
		{
			name: "less",
			v1:   valUint(334344455, expr.Width32),
			v2:   valUint(3874344455, expr.Width32),
			w:    expr.Width32,
			ltu:  true,
			leu:  true,
		},
		{
			name: "greater",
			v1:   valUint(54356, expr.Width16),
			v2:   valUint(33456, expr.Width16),
			w:    expr.Width16,
			ltu:  false,
			leu:  false,
		},
		{
			name: "cut_extend",
			v1:   value{45, 135, 0, 0, 34, 67, 87},
			v2:   value{45, 135},
			w:    expr.Width32,
			ltu:  false,
			leu:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("ltu", func(t *testing.T) {
				res := ltu(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.ltu, res)
			})

			t.Run("leu", func(t *testing.T) {
				res := leu(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.leu, res)
			})
		})
	}
}

func TestLtsAndLes(t *testing.T) {
	tests := []struct {
		name string
		v1   value
		v2   value
		w    expr.Width
		lts  bool
		les  bool
	}{
		{
			name: "equal_signed",
			v1:   value{5, 7, 33, 249},
			v2:   value{5, 7, 33, 249},
			w:    expr.Width32,
			lts:  false,
			les:  true,
		},
		{
			name: "equal_unsigned",
			v1:   value{5, 7, 33, 49},
			v2:   value{5, 7, 33, 49},
			w:    expr.Width32,
			lts:  false,
			les:  true,
		},
		{
			name: "less_negative",
			v1:   value{0x4f, 0x23, 0x8c, 0x81},
			v2:   value{0x4f, 0x23, 0x8c, 0x80},
			w:    expr.Width32,
			lts:  true,
			les:  true,
		},
		{
			name: "less_positive",
			v1:   value{0x4f, 0x23, 0x8c, 0x71},
			v2:   value{0x4f, 0x24, 0x8c, 0x71},
			w:    expr.Width32,
			lts:  true,
			les:  true,
		},
		{
			name: "less_negative_positive",
			v1:   value{0x4f, 0x23, 0x8c, 0x81},
			v2:   value{0x4f, 0x23, 0x8c, 0x76},
			w:    expr.Width32,
			lts:  true,
			les:  true,
		},
		{
			name: "cut_extend",
			v1:   value{45, 135, 0, 0, 34, 67, 87},
			v2:   value{45, 135},
			w:    expr.Width32,
			lts:  false,
			les:  true,
		},
		{
			name: "zero_and_minus_one",
			v1:   value{0},
			v2:   value{0xff, 0xff, 0xff, 0xff},
			w:    expr.Width32,
			lts:  false,
			les:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("lts", func(t *testing.T) {
				res := lts(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.lts, res)
			})

			t.Run("les", func(t *testing.T) {
				res := les(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.les, res)
			})
		})
	}
}
