package expreval_test

import (
	"mltwist/internal/exprtransform/internal/expreval"
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEq(t *testing.T) {
	tests := []struct {
		name     string
		v1       expreval.Value
		v2       expreval.Value
		w        expr.Width
		expected bool
	}{
		{
			name:     "equal",
			v1:       expreval.NewValue([]byte{5, 7, 33, 249}),
			v2:       expreval.NewValue([]byte{5, 7, 33, 249}),
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
			v1:       expreval.NewValue([]byte{45, 135, 0, 0, 34, 67, 87}),
			v2:       expreval.NewValue([]byte{45, 135}),
			w:        expr.Width32,
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := expreval.Eq(tt.v1, tt.v2, tt.w)
			require.Equal(t, tt.expected, res)
		})
	}
}

func TestLtu(t *testing.T) {
	tests := []struct {
		name string
		v1   expreval.Value
		v2   expreval.Value
		w    expr.Width
		ltu  bool
	}{
		{
			name: "equal",
			v1:   expreval.NewValue([]byte{5, 7, 33, 249}),
			v2:   expreval.NewValue([]byte{5, 7, 33, 249}),
			w:    expr.Width32,
			ltu:  false,
		},
		{
			name: "less",
			v1:   valUint(334344455, expr.Width32),
			v2:   valUint(3874344455, expr.Width32),
			w:    expr.Width32,
			ltu:  true,
		},
		{
			name: "greater",
			v1:   valUint(54356, expr.Width16),
			v2:   valUint(33456, expr.Width16),
			w:    expr.Width16,
			ltu:  false,
		},
		{
			name: "cut_extend",
			v1:   expreval.NewValue([]byte{45, 135, 0, 0, 34, 67, 87}),
			v2:   expreval.NewValue([]byte{45, 135}),
			w:    expr.Width32,
			ltu:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := expreval.Ltu(tt.v1, tt.v2, tt.w)
			require.Equal(t, tt.ltu, res)
		})
	}
}

func TestLts(t *testing.T) {
	tests := []struct {
		name string
		v1   expreval.Value
		v2   expreval.Value
		w    expr.Width
		lts  bool
	}{
		{
			name: "equal_signed",
			v1:   expreval.NewValue([]byte{5, 7, 33, 249}),
			v2:   expreval.NewValue([]byte{5, 7, 33, 249}),
			w:    expr.Width32,
			lts:  false,
		},
		{
			name: "equal_unsigned",
			v1:   expreval.NewValue([]byte{5, 7, 33, 49}),
			v2:   expreval.NewValue([]byte{5, 7, 33, 49}),
			w:    expr.Width32,
			lts:  false,
		},
		{
			name: "less_negative",
			v1:   expreval.NewValue([]byte{0x4f, 0x23, 0x8c, 0x81}),
			v2:   expreval.NewValue([]byte{0x4f, 0x23, 0x8c, 0x80}),
			w:    expr.Width32,
			lts:  true,
		},
		{
			name: "less_positive",
			v1:   expreval.NewValue([]byte{0x4f, 0x23, 0x8c, 0x71}),
			v2:   expreval.NewValue([]byte{0x4f, 0x24, 0x8c, 0x71}),
			w:    expr.Width32,
			lts:  true,
		},
		{
			name: "less_negative_positive",
			v1:   expreval.NewValue([]byte{0x4f, 0x23, 0x8c, 0x81}),
			v2:   expreval.NewValue([]byte{0x4f, 0x23, 0x8c, 0x76}),
			w:    expr.Width32,
			lts:  true,
		},
		{
			name: "cut_extend",
			v1:   expreval.NewValue([]byte{45, 135, 0, 0, 34, 67, 87}),
			v2:   expreval.NewValue([]byte{45, 135}),
			w:    expr.Width32,
			lts:  false,
		},
		{
			name: "zero_and_minus_one",
			v1:   expreval.NewValue([]byte{0}),
			v2:   expreval.NewValue([]byte{0xff, 0xff, 0xff, 0xff}),
			w:    expr.Width32,
			lts:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := expreval.Lts(tt.v1, tt.v2, tt.w)
			require.Equal(t, tt.lts, res)
		})
	}
}
