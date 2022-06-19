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
		cnt  exprtools.BitCnt
		w    expr.Width
		exp  expr.Expr
	}{
		{
			name: "singlebyte",
			e:    expr.ConstFromUint[uint8](0x7c),
			cnt:  5,
			w:    expr.Width16,
			exp:  expr.ConstFromUint[uint16](0x1c),
		},
		{
			name: "multibyte",
			e:    expr.ConstFromUint[uint16](0xffff),
			cnt:  12,
			w:    expr.Width16,
			exp:  expr.ConstFromUint[uint16](0x0fff),
		},
		{
			name: "bit_int",
			e:    expr.NewConst([]byte{1, 2, 3, 4, 5, 6, 7, 8, 0xca, 0xff, 0xff, 0xff}, expr.Width128),
			cnt:  82,
			w:    expr.Width128,
			exp:  expr.NewConst([]byte{1, 2, 3, 4, 5, 6, 7, 8, 0xca, 0xff, 0x3}, expr.Width128),
		},
		{
			name: "uint64_width",
			e:    expr.ConstFromUint[uint64](0xfedcba9876543210),
			cnt:  64,
			w:    expr.Width64,
			exp:  expr.ConstFromUint[uint64](0xfedcba9876543210),
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
		e:    expr.ConstFromInt[int16](math.MaxInt16),
		w:    expr.Width16,
		exp:  expr.ConstFromUint[uint16](0),
	}, {
		name: "negative",
		e:    expr.ConstFromInt[int32](-1),
		w:    expr.Width32,
		exp:  expr.ConstFromUint[uint32](0x80000000),
	}, {
		name: "grown_to_positive",
		e:    expr.ConstFromInt[int32](-1),
		w:    expr.Width64,
		exp:  expr.ConstFromUint[uint64](0),
	}, {
		name: "cropped_to_negative",
		e:    expr.ConstFromUint[uint16](0x10ff),
		w:    expr.Width8,
		exp:  expr.ConstFromUint[uint8](0x80),
	}, {
		name: "zero",
		e:    expr.ConstFromUint[uint16](0),
		w:    expr.Width16,
		exp:  expr.ConstFromUint[uint16](0),
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
