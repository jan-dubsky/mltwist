package exprtools_test

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitNot(t *testing.T) {
	tests := []struct {
		name string
		e    expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{{
		name: "keep_width",
		e:    expr.ConstFromUint[uint64](0xfedbca9876543210),
		w:    expr.Width64,
		exp:  expr.ConstFromUint(^uint64(0xfedbca9876543210)),
	}, {
		name: "cut_width",
		e:    expr.ConstFromUint[uint64](0xfedbca9876543210),
		w:    expr.Width16,
		exp:  expr.ConstFromUint(^uint16(0x3210)),
	}, {
		name: "grow_width",
		e:    expr.ConstFromUint[uint16](0xb3d7),
		w:    expr.Width32,
		exp:  expr.ConstFromUint(^uint32(0xb3d7)),
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.BitNot(tt.e, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestAndOrXor(t *testing.T) {
	tests := []struct {
		name string
		e1   expr.Expr
		e2   expr.Expr
		w    expr.Width

		and expr.Expr
		or  expr.Expr
		xor expr.Expr
	}{{
		name: "keep_width",
		e1:   expr.ConstFromUint[uint64](0xfedbca9876543210),
		e2:   expr.ConstFromUint[uint64](0x0123456789abcdef),
		w:    expr.Width64,
		and:  expr.ConstFromUint[uint64](0xfedbca9876543210 & 0x0123456789abcdef),
		or:   expr.ConstFromUint[uint64](0xfedbca9876543210 | 0x0123456789abcdef),
		xor:  expr.ConstFromUint[uint64](0xfedbca9876543210 ^ 0x0123456789abcdef),
	}, {
		name: "cut_width",
		e1:   expr.ConstFromUint[uint64](0x23b273ad730509aa),
		e2:   expr.ConstFromUint[uint64](0x9232badfe2913ad3),
		w:    expr.Width32,
		and:  expr.ConstFromUint[uint32](0x730509aa & 0xe2913ad3),
		or:   expr.ConstFromUint[uint32](0x730509aa | 0xe2913ad3),
		xor:  expr.ConstFromUint[uint32](0x730509aa ^ 0xe2913ad3),
	}, {
		name: "gro_width",
		e1:   expr.ConstFromUint[uint8](0x35),
		e2:   expr.ConstFromUint[uint16](0x94b2),
		w:    expr.Width32,
		and:  expr.ConstFromUint[uint32](0x35 & 0x94b2),
		or:   expr.ConstFromUint[uint32](0x35 | 0x94b2),
		xor:  expr.ConstFromUint[uint32](0x35 ^ 0x94b2),
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("and", func(t *testing.T) {
				e := exprtools.BitAnd(tt.e1, tt.e2, tt.w)
				require.Equal(t, tt.and, exprtransform.ConstFold(e))
			})
			t.Run("or", func(t *testing.T) {
				e := exprtools.BitOr(tt.e1, tt.e2, tt.w)
				require.Equal(t, tt.or, exprtransform.ConstFold(e))
			})
			t.Run("xor", func(t *testing.T) {
				e := exprtools.BitXor(tt.e1, tt.e2, tt.w)
				require.Equal(t, tt.xor, exprtransform.ConstFold(e))
			})
		})
	}
}
