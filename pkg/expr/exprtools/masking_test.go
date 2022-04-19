package exprtools_test

import (
	"math"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMaskBits(t *testing.T) {
	tests := []struct {
		name string
		e    expr.Expr
		cnt  uint16
		w    expr.Width
		exp  expr.Expr
	}{
		{
			name: "singlebyte",
			e:    expr.NewConstUint[uint8](0x7c, expr.Width8),
			cnt:  5,
			w:    expr.Width16,
			exp:  expr.NewConstUint[uint16](0x1c, expr.Width16),
		},
		{
			name: "multibyte",
			e:    expr.NewConstUint[uint16](0xffff, expr.Width16),
			cnt:  12,
			w:    expr.Width16,
			exp:  expr.NewConstUint[uint16](0x0fff, expr.Width16),
		},
		{
			name: "bit_int",
			e:    expr.NewConst([]byte{1, 2, 3, 4, 5, 6, 7, 8, 0xca, 0xff, 0xff, 0xff}, expr.Width128),
			cnt:  82,
			w:    expr.Width128,
			exp:  expr.NewConst([]byte{1, 2, 3, 4, 5, 6, 7, 8, 0xca, 0xff, 0x3}, expr.Width128),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.MaskBits(tt.e, tt.cnt, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestIntNegative(t *testing.T) {
	tests := []struct {
		name string
		e    expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{{
		name: "positive",
		e:    expr.NewConstInt[int16](math.MaxInt16, expr.Width16),
		w:    expr.Width16,
		exp:  expr.NewConstUint[uint8](0, expr.Width16),
	}, {
		name: "negative",
		e:    expr.NewConstInt[int32](-1, expr.Width32),
		w:    expr.Width32,
		exp:  expr.NewConstUint[uint32](0x80000000, expr.Width32),
	}, {
		name: "grown_to_positive",
		e:    expr.NewConstInt[int32](-1, expr.Width32),
		w:    expr.Width64,
		exp:  expr.NewConstUint[uint8](0, expr.Width64),
	}, {
		name: "cropped_to_negative",
		e:    expr.NewConstUint[uint16](0x10ff, expr.Width16),
		w:    expr.Width8,
		exp:  expr.NewConstUint[uint8](0x80, expr.Width8),
	}, {
		name: "zero",
		e:    expr.NewConstUint[uint16](0, expr.Width16),
		w:    expr.Width16,
		exp:  expr.NewConstUint[uint16](0, expr.Width16),
	}, {
		name: "big_positive",
		e:    expr.NewConst([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4}, expr.Width128),
		w:    expr.Width128,
		exp:  expr.NewConstUint[uint8](0, expr.Width128),
	}, {
		name: "big_negative",
		e:    expr.NewConst([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 128 + 4}, expr.Width128),
		w:    expr.Width128,
		exp:  expr.NewConst([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128}, expr.Width128),
	},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.IntNegative(tt.e, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}
