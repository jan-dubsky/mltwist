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
			e:    expr.ConstFromUint[uint8](57),
			w:    expr.Width32,
			exp:  expr.ConstFromInt[int32](-57),
		},
		{
			name: "neg_to_pos",
			e:    expr.ConstFromInt[int16](-255),
			w:    expr.Width8,
			exp:  expr.ConstFromUint[uint8](255),
		},
		{
			name: "zero",
			e:    expr.ConstFromInt[int16](0),
			w:    expr.Width16,
			exp:  expr.ConstFromUint[uint16](0),
		},
		{
			name: "pos_to_neg_multibyte",
			e:    expr.ConstFromUint[uint16](0x115e),
			w:    expr.Width64,
			exp:  expr.ConstFromInt[int64](-0x115e),
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
			e:    expr.ConstFromUint[uint8](57),
			w:    expr.Width8,
			exp:  expr.ConstFromUint[uint8](57),
		},
		{
			name: "positive_width_change",
			e:    expr.ConstFromUint[uint8](57),
			w:    expr.Width32,
			exp:  expr.ConstFromUint[uint32](57),
		},
		{
			name: "negative",
			e:    expr.ConstFromInt[int16](-0x5c3b),
			w:    expr.Width16,
			exp:  expr.ConstFromUint[uint16](0x5c3b),
		},
		{
			name: "zero",
			e:    expr.ConstFromInt[int16](0),
			w:    expr.Width16,
			exp:  expr.ConstFromUint[uint16](0),
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
			exp: expr.ConstFromUint[uint16](0xffff),
		},
		{
			w:   expr.Width8,
			exp: expr.ConstFromUint[uint8](0xff),
		},
		{
			w:   expr.Width64,
			exp: expr.ConstFromUint[uint64](0xffffffffffffffff),
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

func TestMod(t *testing.T) {
	tests := []struct {
		name string
		e1   expr.Expr
		e2   expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{{
		name: "plain_modulo",
		e1:   expr.ConstFromUint[uint16](46323),
		e2:   expr.ConstFromUint[uint8](134),
		w:    expr.Width16,
		exp:  expr.NewConstUint[uint16](46323%134, expr.Width16),
	}, {
		name: "mod_by_zero",
		e1:   expr.ConstFromUint[uint16](46323),
		e2:   expr.Zero,
		w:    expr.Width32,
		exp:  expr.NewConstUint[uint16](46323, expr.Width32),
	}, {
		name: "mod_by_one",
		e1:   expr.ConstFromInt[int32](-1),
		e2:   expr.One,
		w:    expr.Width32,
		exp:  expr.NewConstUint[uint8](0, expr.Width32),
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := exprtools.Mod(tt.e1, tt.e2, tt.w)
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
			e1:   expr.ConstFromUint[uint8](57),
			e2:   expr.ConstFromUint[uint8](37),
			w:    expr.Width8,
			exp:  expr.ConstFromUint[uint16](57 * 37),
		},
		{
			name: "pos_mul_neg",
			e1:   expr.ConstFromInt[int8](57),
			e2:   expr.ConstFromInt[int32](-78),
			w:    expr.Width32,
			exp:  expr.ConstFromInt[int64](-(57 * 78)),
		},
		{
			name: "neg_mul_pos",
			e1:   expr.ConstFromInt[int16](-0x5c3b),
			e2:   expr.ConstFromUint[uint8](2),
			w:    expr.Width16,
			exp:  expr.ConstFromInt[int32](-(0x5c3b * 2)),
		},
		{
			name: "neg_mul_neg",
			e1:   expr.ConstFromInt[int32](-7),
			e2:   expr.ConstFromInt[int64](-2),
			w:    expr.Width8,
			exp:  expr.ConstFromUint[uint16](14),
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
				expr.ConstFromInt[int32](-7),
				expr.ConstFromInt[int64](-2),
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
			e1:     expr.ConstFromUint[uint8](57),
			e2:     expr.ConstFromUint[uint8](37),
			w:      expr.Width16,
			expDiv: expr.ConstFromUint[uint16](1),
			expMod: expr.ConstFromUint[uint16](57 % 37),
		}, {
			name:   "pos_by_neg_lt",
			e1:     expr.ConstFromInt[int8](57),
			e2:     expr.ConstFromInt[int32](-78),
			w:      expr.Width32,
			expDiv: expr.ConstFromInt[int32](0),
			expMod: expr.ConstFromInt[int32](-57),
		},
		{
			name:   "pos_by_neg",
			e1:     expr.ConstFromInt[int16](998),
			e2:     expr.ConstFromInt[int32](-78),
			w:      expr.Width64,
			expDiv: expr.ConstFromInt[int64](-998 / 78),
			expMod: expr.ConstFromInt[int64](-998 % 78),
		},
		{
			name:   "neg_by_pos",
			e1:     expr.ConstFromInt[int16](-0x5c3b),
			e2:     expr.ConstFromUint[uint8](2),
			w:      expr.Width16,
			expDiv: expr.ConstFromInt[int16](-(0x5c3b / 2)),
			expMod: expr.ConstFromInt[int16](-1),
		},
		{
			name:   "neg_by_neg",
			e1:     expr.ConstFromInt[int32](-74),
			e2:     expr.ConstFromInt[int64](-23),
			w:      expr.Width8,
			expDiv: expr.ConstFromUint[uint8](74 / 23),
			expMod: expr.ConstFromUint[uint8](74 % 23),
		},
		{
			name:   "pos_by_zero",
			e1:     expr.ConstFromInt[int32](0x1ccd30d8),
			e2:     expr.Zero,
			w:      expr.Width16,
			expDiv: expr.ConstFromUint[uint16](0xffff),
			expMod: expr.ConstFromUint[uint16](0x30d8),
		},
		{
			name:   "neg_by_zero",
			e1:     expr.ConstFromInt[int64](-0xffbc),
			e2:     expr.Zero,
			w:      expr.Width64,
			expDiv: expr.ConstFromInt[int64](-1),
			expMod: expr.ConstFromInt[int64](-0xffbc),
		},
		{
			name:   "overflow",
			e1:     expr.ConstFromInt[int32](math.MinInt32),
			e2:     expr.ConstFromInt[int32](-1),
			w:      expr.Width32,
			expDiv: expr.ConstFromInt[int32](math.MinInt32),
			expMod: expr.ConstFromInt[int32](0),
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
			e:    expr.ConstFromInt[int8](127),
			sign: expr.ConstFromUint[uint8](7),
			w:    expr.Width16,
			exp:  expr.ConstFromInt[int16](127),
		},
		{
			name: "extend_negative",
			e:    expr.ConstFromInt[int16](-567),
			sign: expr.ConstFromUint[uint8](15),
			w:    expr.Width32,
			exp:  expr.ConstFromInt[int32](-567),
		},
		{
			name: "positive_higher_bits_set",
			e:    expr.ConstFromUint[uint16](0x1234),
			sign: expr.ConstFromUint[uint8](7),
			w:    expr.Width32,
			exp:  expr.ConstFromInt[int32](0x34),
		},
		{
			name: "negative_higher_bits_unset",
			e:    expr.ConstFromUint[uint16](0x00b4),
			sign: expr.ConstFromUint[uint8](7),
			w:    expr.Width32,
			exp:  expr.ConstFromUint[uint32](0xffffffb4),
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
			e:     expr.ConstFromUint[uint16](0x7b5a),
			shift: expr.ConstFromUint[uint8](5),
			w:     expr.Width16,
			exp:   expr.ConstFromUint[uint16](0x7b5a >> 5),
		},
		{
			name:  "negative",
			e:     expr.ConstFromUint[uint16](0xbb5a),
			shift: expr.ConstFromUint[uint8](3),
			w:     expr.Width16,
			exp:   expr.ConstFromUint[uint16]((0xbb5a >> 3) | 0xe000),
		},
		{
			name:  "positive_shifted_many",
			e:     expr.ConstFromInt[int32](0x593290ba),
			shift: expr.ConstFromUint[uint8](236),
			w:     expr.Width64,
			exp:   expr.ConstFromUint[uint64](0),
		},
		{
			name:  "negative_shifted_many",
			e:     expr.ConstFromInt[int32](-2143044594),
			shift: expr.ConstFromUint[uint8](156),
			w:     expr.Width32,
			exp:   expr.ConstFromInt[int32](-1),
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
