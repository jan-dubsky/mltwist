package expreval_test

import (
	"mltwist/internal/exprtransform/internal/expreval"
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

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
