package exprtools_test

import (
	"fmt"
	"math"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNegate(t *testing.T) {
	tests := []struct {
		name string
		e    expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{
		{
			name: "pos_to_neg",
			e:    expr.NewConstUint[uint8](57, expr.Width8),
			w:    expr.Width32,
			exp:  expr.NewConstInt[int8](-57, expr.Width32),
		},
		{
			name: "neg_to_pos",
			e:    expr.NewConstInt[int16](-255, expr.Width16),
			w:    expr.Width8,
			exp:  expr.NewConstUint[uint8](255, expr.Width8),
		},
		{
			name: "zero",
			e:    expr.NewConstInt[int16](0, expr.Width16),
			w:    expr.Width16,
			exp:  expr.NewConstUint[uint16](0, expr.Width16),
		},
		{
			name: "pos_to_neg_multibyte",
			e:    expr.NewConstUint[uint16](0x115e, expr.Width16),
			w:    expr.Width64,
			exp:  expr.NewConstInt[int16](-0x115e, expr.Width64),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.Negate(tt.e, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		e    expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{
		{
			name: "positive",
			e:    expr.NewConstUint[uint8](57, expr.Width8),
			w:    expr.Width8,
			exp:  expr.NewConstUint[uint8](57, expr.Width8),
		},
		{
			name: "positive_width_change",
			e:    expr.NewConstUint[uint8](57, expr.Width8),
			w:    expr.Width32,
			exp:  expr.NewConstUint[uint8](57, expr.Width32),
		},
		{
			name: "negative",
			e:    expr.NewConstInt[int16](-0x5c3b, expr.Width16),
			w:    expr.Width16,
			exp:  expr.NewConstUint[uint16](0x5c3b, expr.Width16),
		},
		{
			name: "zero",
			e:    expr.NewConstInt[int16](0, expr.Width16),
			w:    expr.Width16,
			exp:  expr.NewConstUint[uint16](0, expr.Width16),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.Abs(tt.e, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestOnes(t *testing.T) {
	tests := []struct {
		w   expr.Width
		exp expr.Expr
	}{
		{
			w:   expr.Width16,
			exp: expr.NewConstUint[uint16](0xffff, expr.Width16),
		},
		{
			w:   expr.Width8,
			exp: expr.NewConstUint[uint8](0xff, expr.Width8),
		},
		{
			w:   expr.Width64,
			exp: expr.NewConstUint[uint64](0xffffffffffffffff, expr.Width64),
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			e := exprtools.Ones(tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestSignedMul(t *testing.T) {
	tests := []struct {
		name string
		e1   expr.Expr
		e2   expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{
		{
			name: "pos_mul_pos",
			e1:   expr.NewConstUint[uint8](57, expr.Width8),
			e2:   expr.NewConstUint[uint8](37, expr.Width8),
			w:    expr.Width8,
			exp:  expr.NewConstUint[uint16](57*37, expr.Width16),
		},
		{
			name: "pos_mul_neg",
			e1:   expr.NewConstInt[int8](57, expr.Width8),
			e2:   expr.NewConstInt[int8](-78, expr.Width32),
			w:    expr.Width32,
			exp:  expr.NewConstInt[int16](-(57 * 78), expr.Width64),
		},
		{
			name: "neg_mul_pos",
			e1:   expr.NewConstInt[int16](-0x5c3b, expr.Width16),
			e2:   expr.NewConstUint[uint8](2, expr.Width8),
			w:    expr.Width16,
			exp:  expr.NewConstInt[int32](-(0x5c3b * 2), expr.Width32),
		},
		{
			name: "neg_mul_neg",
			e1:   expr.NewConstInt[int16](-7, expr.Width32),
			e2:   expr.NewConstInt[int16](-2, expr.Width64),
			w:    expr.Width8,
			exp:  expr.NewConstUint[uint8](14, expr.Width16),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.SignedMul(tt.e1, tt.e2, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}

	t.Run("too_wide", func(t *testing.T) {
		require.Panics(t, func() {
			exprtools.SignedMul(
				expr.NewConstInt[int16](-7, expr.Width32),
				expr.NewConstInt[int16](-2, expr.Width64),
				expr.MaxWidth,
			)
		})
	})
}

func TestSignedDivMod(t *testing.T) {
	tests := []struct {
		name   string
		e1     expr.Expr
		e2     expr.Expr
		w      expr.Width
		expDiv expr.Expr
		expMod expr.Expr
	}{
		{
			name:   "pos_by_pos",
			e1:     expr.NewConstUint[uint8](57, expr.Width8),
			e2:     expr.NewConstUint[uint8](37, expr.Width8),
			w:      expr.Width16,
			expDiv: expr.NewConstUint[uint16](1, expr.Width16),
			expMod: expr.NewConstUint[uint16](57%37, expr.Width16),
		}, {
			name:   "pos_by_neg_lt",
			e1:     expr.NewConstInt[int8](57, expr.Width8),
			e2:     expr.NewConstInt[int8](-78, expr.Width32),
			w:      expr.Width32,
			expDiv: expr.NewConstInt[int16](0, expr.Width32),
			expMod: expr.NewConstInt[int16](-57, expr.Width32),
		},
		{
			name:   "pos_by_neg",
			e1:     expr.NewConstInt[int16](998, expr.Width16),
			e2:     expr.NewConstInt[int8](-78, expr.Width32),
			w:      expr.Width64,
			expDiv: expr.NewConstInt[int16](-998/78, expr.Width64),
			expMod: expr.NewConstInt[int16](-998%78, expr.Width64),
		},
		{
			name:   "neg_by_pos",
			e1:     expr.NewConstInt[int16](-0x5c3b, expr.Width16),
			e2:     expr.NewConstUint[uint8](2, expr.Width8),
			w:      expr.Width16,
			expDiv: expr.NewConstInt[int32](-(0x5c3b / 2), expr.Width16),
			expMod: expr.NewConstInt[int32](-1, expr.Width16),
		},
		{
			name:   "neg_by_neg",
			e1:     expr.NewConstInt[int16](-74, expr.Width32),
			e2:     expr.NewConstInt[int16](-23, expr.Width64),
			w:      expr.Width8,
			expDiv: expr.NewConstUint[uint8](74/23, expr.Width8),
			expMod: expr.NewConstUint[uint8](74%23, expr.Width8),
		},
		{
			name:   "pos_by_zero",
			e1:     expr.NewConstInt[int32](0x1ccd30d8, expr.Width32),
			e2:     expr.Zero,
			w:      expr.Width16,
			expDiv: expr.NewConstUint[uint16](0xffff, expr.Width16),
			expMod: expr.NewConstUint[uint16](0x30d8, expr.Width16),
		},
		{
			name:   "neg_by_zero",
			e1:     expr.NewConstInt[int64](-0xffbc, expr.Width64),
			e2:     expr.Zero,
			w:      expr.Width64,
			expDiv: expr.NewConstInt[int64](-1, expr.Width64),
			expMod: expr.NewConstInt[int64](-0xffbc, expr.Width64),
		},
		{
			name:   "overflow",
			e1:     expr.NewConstInt[int32](math.MinInt32, expr.Width32),
			e2:     expr.NewConstInt[int32](-1, expr.Width32),
			w:      expr.Width32,
			expDiv: expr.NewConstInt[int32](math.MinInt32, expr.Width32),
			expMod: expr.NewConstInt[int32](0, expr.Width32),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("div", func(t *testing.T) {
				e := exprtools.SignedDiv(tt.e1, tt.e2, tt.w)
				require.Equal(t, tt.expDiv, exprtransform.ConstFold(e))
			})
			t.Run("mod", func(t *testing.T) {
				e := exprtools.SignedMod(tt.e1, tt.e2, tt.w)
				require.Equal(t, tt.expMod, exprtransform.ConstFold(e))
			})
		})
	}
}

func TestSignExtend(t *testing.T) {
	tests := []struct {
		name string
		e    expr.Expr
		sign expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{
		{
			name: "extend_positive",
			e:    expr.NewConstInt[int8](127, expr.Width8),
			sign: expr.NewConstUint[uint8](7, expr.Width8),
			w:    expr.Width16,
			exp:  expr.NewConstInt[int8](127, expr.Width16),
		},
		{
			name: "extend_negative",
			e:    expr.NewConstInt[int16](-567, expr.Width16),
			sign: expr.NewConstUint[uint8](15, expr.Width8),
			w:    expr.Width32,
			exp:  expr.NewConstInt[int16](-567, expr.Width32),
		},
		{
			name: "positive_higher_bits_set",
			e:    expr.NewConstUint[uint16](0x1234, expr.Width16),
			sign: expr.NewConstUint[uint8](7, expr.Width8),
			w:    expr.Width32,
			exp:  expr.NewConstInt[int16](0x34, expr.Width32),
		},
		{
			name: "negative_higher_bits_unset",
			e:    expr.NewConstUint[uint16](0x00b4, expr.Width16),
			sign: expr.NewConstUint[uint8](7, expr.Width8),
			w:    expr.Width32,
			exp:  expr.NewConstUint[uint32](0xffffffb4, expr.Width32),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.SignExtend(tt.e, tt.sign, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestRshA(t *testing.T) {
	tests := []struct {
		name  string
		e     expr.Expr
		shift expr.Expr
		w     expr.Width
		exp   expr.Expr
	}{
		{
			name:  "positive",
			e:     expr.NewConstUint[uint16](0x7b5a, expr.Width16),
			shift: expr.NewConstUint[uint16](5, expr.Width8),
			w:     expr.Width16,
			exp:   expr.NewConstUint[uint16](0x7b5a>>5, expr.Width16),
		},
		{
			name:  "negative",
			e:     expr.NewConstUint[uint16](0xbb5a, expr.Width16),
			shift: expr.NewConstUint[uint16](3, expr.Width8),
			w:     expr.Width16,
			exp:   expr.NewConstUint[uint16]((0xbb5a>>3)|0xe000, expr.Width16),
		},
		{
			name:  "positive_shifted_many",
			e:     expr.NewConstInt[int32](0x593290ba, expr.Width32),
			shift: expr.NewConstUint[uint8](236, expr.Width8),
			w:     expr.Width64,
			exp:   expr.NewConstUint[uint8](0, expr.Width64),
		},
		{
			name:  "negative_shifted_many",
			e:     expr.NewConstInt[int32](-2143044594, expr.Width32),
			shift: expr.NewConstUint[uint8](156, expr.Width8),
			w:     expr.Width32,
			exp:   expr.NewConstInt[int32](-1, expr.Width32),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.RshA(tt.e, tt.shift, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}
