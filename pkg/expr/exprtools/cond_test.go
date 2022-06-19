package exprtools_test

import (
	"math"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEq(t *testing.T) {
	tests := []struct {
		name   string
		arg1   expr.Expr
		arg2   expr.Expr
		w      expr.Width
		isTrue bool
	}{{
		name:   "zero_to_zero",
		arg1:   expr.Zero,
		arg2:   expr.Zero,
		w:      expr.Width16,
		isTrue: true,
	}, {
		name:   "one_to_one",
		arg1:   expr.One,
		arg2:   expr.One,
		w:      expr.Width16,
		isTrue: true,
	}, {
		name:   "zero_to_one",
		arg1:   expr.Zero,
		arg2:   expr.One,
		w:      expr.Width32,
		isTrue: false,
	}, {
		name:   "one_to_zero",
		arg1:   expr.One,
		arg2:   expr.Zero,
		w:      expr.Width32,
		isTrue: false,
	}, {
		name:   "one_to_max",
		arg1:   expr.One,
		arg2:   exprtools.Ones(expr.Width16),
		w:      expr.Width16,
		isTrue: false,
	}, {
		name:   "max_to_max",
		arg1:   exprtools.Ones(expr.Width32),
		arg2:   exprtools.Ones(expr.Width32),
		w:      expr.Width32,
		isTrue: true,
	}, {
		name:   "max_to_max_minus_one",
		arg1:   exprtools.Ones(expr.Width32),
		arg2:   expr.ConstFromUint[uint32](math.MaxUint32 - 1),
		w:      expr.Width32,
		isTrue: false,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			te := expr.ConstFromUint[uint64](0x0123456789abcdef)
			fe := expr.ConstFromUint[uint64](0x84848484bcbcbcbc)

			e := exprtools.Eq(tt.arg1, tt.arg2, te, fe, tt.w)

			exp := fe
			if tt.isTrue {
				exp = te
			}
			exp = exp.WithWidth(tt.w)

			require.Equal(t, exp, exprtransform.ConstFold(e))
		})
	}
}

func TestLts(t *testing.T) {
	tests := []struct {
		name   string
		arg1   expr.Expr
		arg2   expr.Expr
		w      expr.Width
		isTrue bool
	}{{
		name:   "zero_to_zero",
		arg1:   expr.Zero,
		arg2:   expr.Zero,
		w:      expr.Width16,
		isTrue: false,
	}, {
		name:   "zero_to_one",
		arg1:   expr.Zero,
		arg2:   expr.One,
		w:      expr.Width32,
		isTrue: true,
	}, {
		name:   "one_to_zero",
		arg1:   expr.One,
		arg2:   expr.Zero,
		w:      expr.Width32,
		isTrue: false,
	}, {
		name:   "one_to_minus_max",
		arg1:   expr.One,
		arg2:   exprtools.Ones(expr.Width16),
		w:      expr.Width16,
		isTrue: false,
	}, {
		name:   "minus_max_to_one",
		arg1:   exprtools.Ones(expr.Width16),
		arg2:   expr.One,
		w:      expr.Width16,
		isTrue: true,
	}, {
		name:   "minus_max_to_minus_max",
		arg1:   exprtools.Ones(expr.Width32),
		arg2:   exprtools.Ones(expr.Width32),
		w:      expr.Width32,
		isTrue: false,
	}, {
		name:   "zero_to_minus_one",
		arg1:   expr.Zero,
		arg2:   expr.ConstFromInt[int16](-1),
		w:      expr.Width16,
		isTrue: false,
	}, {
		name:   "minus_one_to_zero",
		arg1:   expr.ConstFromInt[int16](-1),
		arg2:   expr.Zero,
		w:      expr.Width16,
		isTrue: true,
	}, {
		name:   "minus_one_to_minus_five",
		arg1:   expr.ConstFromInt[int16](-1),
		arg2:   expr.ConstFromInt[int16](-5),
		w:      expr.Width16,
		isTrue: false,
	}, {
		name:   "minus_five_to_minus_one",
		arg1:   expr.ConstFromInt[int16](-5),
		arg2:   expr.ConstFromInt[int16](-1),
		w:      expr.Width16,
		isTrue: true,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			te := expr.ConstFromUint[uint64](0x0123456789abcdef)
			fe := expr.ConstFromUint[uint64](0x84848484bcbcbcbc)

			e := exprtools.Lts(tt.arg1, tt.arg2, te, fe, tt.w)

			exp := fe
			if tt.isTrue {
				exp = te
			}
			exp = exp.WithWidth(tt.w)

			require.Equal(t, exp, exprtransform.ConstFold(e))
		})
	}
}

func TestLeu(t *testing.T) {
	tests := []struct {
		name   string
		arg1   expr.Expr
		arg2   expr.Expr
		w      expr.Width
		isTrue bool
	}{{
		name:   "one_to_zero",
		arg1:   expr.Zero,
		arg2:   expr.One,
		w:      expr.Width32,
		isTrue: true,
	}, {
		name:   "one_to_zero",
		arg1:   expr.Zero,
		arg2:   expr.Zero,
		w:      expr.Width16,
		isTrue: true,
	}, {
		name: "one_to_zero_wide",
		arg1: expr.One,
		arg2: expr.Zero,
		w:    expr.Width64,
	}, {
		name:   "value_to_minus_one",
		arg1:   expr.ConstFromUint[uint8](53),
		arg2:   exprtools.Ones(expr.Width32),
		w:      expr.Width32,
		isTrue: true,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			te := expr.ConstFromUint[uint64](0x0123456789abcdef)
			fe := expr.ConstFromUint[uint64](0x84848484bcbcbcbc)

			e := exprtools.Leu(tt.arg1, tt.arg2, te, fe, tt.w)

			exp := fe
			if tt.isTrue {
				exp = te
			}
			exp = exp.WithWidth(tt.w)

			require.Equal(t, exp, exprtransform.ConstFold(e))
		})
	}
}

func TestLes(t *testing.T) {
	tests := []struct {
		name   string
		arg1   expr.Expr
		arg2   expr.Expr
		w      expr.Width
		isTrue bool
	}{{
		name:   "one_to_zero",
		arg1:   expr.Zero,
		arg2:   expr.One,
		w:      expr.Width32,
		isTrue: true,
	}, {
		name:   "one_to_zero",
		arg1:   expr.Zero,
		arg2:   expr.Zero,
		w:      expr.Width16,
		isTrue: true,
	}, {
		name: "one_to_zero",
		arg1: expr.One,
		arg2: expr.Zero,
		w:    expr.Width64,
	}, {
		name: "value_to_minus_one",
		arg1: expr.ConstFromUint[uint8](53),
		arg2: exprtools.Ones(expr.Width32),
		w:    expr.Width32,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			te := expr.ConstFromUint[uint64](0x0123456789abcdef)
			fe := expr.ConstFromUint[uint64](0x84848484bcbcbcbc)

			e := exprtools.Les(tt.arg1, tt.arg2, te, fe, tt.w)

			exp := fe
			if tt.isTrue {
				exp = te
			}
			exp = exp.WithWidth(tt.w)

			require.Equal(t, exp, exprtransform.ConstFold(e))
		})
	}
}
