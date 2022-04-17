package expreval_test

import (
	"mltwist/internal/exprtransform/internal/expreval"
	"mltwist/pkg/expr"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type binaryOpTest struct {
	name   string
	v1     expreval.Value
	v2     expreval.Value
	w      expr.Width
	result expreval.Value
}

type binaryOpFunc func(v1 expreval.Value, v2 expreval.Value, w expr.Width) expreval.Value

func testBinaryOp(t *testing.T, tests []binaryOpTest, f binaryOpFunc) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := f(tt.v1, tt.v2, tt.w)
			require.Equal(t, tt.result, res)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "single_byte",
			v1:     valUint(76, expr.Width8),
			v2:     valUint(32, expr.Width8),
			w:      expr.Width8,
			result: valUint(76+32, expr.Width8),
		},
		{
			name:   "sum_with_carry",
			v1:     expreval.Value{255},
			v2:     expreval.Value{1},
			w:      expr.Width32,
			result: expreval.Value{0, 1, 0, 0},
		},
		{
			name:   "carry_chain",
			v1:     expreval.Value{255, 255, 255, 255, 255},
			v2:     expreval.Value{7, 67, 23, 43},
			w:      expr.Width64,
			result: expreval.Value{6, 67, 23, 43, 0, 1, 0, 0},
		},
		{
			name:   "cut_upper_bytes",
			v1:     expreval.Value{5, 6, 7, 8},
			v2:     expreval.Value{34, 255},
			w:      expr.Width16,
			result: expreval.Value{39, 5},
		},
	}
	testBinaryOp(t, tests, expreval.Add)
}

func TestSub(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "without_carry",
			v1:     expreval.Value{45, 67, 89, 8},
			v2:     expreval.Value{45, 45, 84, 8},
			w:      expr.Width32,
			result: expreval.Value{0, 22, 5, 0},
		},
		{
			name:   "carry_chain",
			v1:     expreval.Value{98, 93, 34, 67, 89, 78},
			v2:     expreval.Value{99, 245, 88, 98, 0, 12},
			w:      expr.Width64,
			result: expreval.Value{255, 103, 201, 224, 88, 66, 0, 0},
		},
		{
			name:   "cut_upper_bytes",
			v1:     expreval.Value{5, 6, 7, 8},
			v2:     expreval.Value{34, 254},
			w:      expr.Width16,
			result: expreval.Value{227, 7},
		},
	}
	testBinaryOp(t, tests, expreval.Sub)
}

func TestLsh(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "byte_shift",
			v1:     expreval.Value{1, 2, 3, 4},
			v2:     valUint(16, expr.Width8),
			w:      expr.Width32,
			result: expreval.Value{0, 0, 1, 2},
		},
		{
			name:   "bit_shift",
			v1:     expreval.Value{0xff, 0, 0x0f, 0x2},
			v2:     valUint(3, expr.Width8),
			w:      expr.Width32,
			result: expreval.Value{0xf8, 7, 0x78, 0x10},
		},
		{
			name:   "shift_cut_to_zero",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     expreval.Value{0, 0, 0, 0, 0, 0, 0, 0, 1},
			w:      expr.Width64,
			result: valUint(math.MaxUint64, expr.Width64),
		},
		{
			name:   "shift_above_uint64",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     expreval.Value{0, 0, 0, 0, 0, 0, 0, 0, 1},
			w:      expr.Width128,
			result: valUint(0, expr.Width128),
		},
		{
			name:   "byte_shift_equals_width",
			v1:     valUint(0xf423b129, expr.Width64),
			v2:     valUint(65, expr.Width8),
			w:      expr.Width64,
			result: valUint(0, expr.Width64),
		},
		{
			name:   "byte_shift_above_width",
			v1:     valUint(0xf423b129, expr.Width64),
			v2:     valUint(9*8+2, expr.Width8),
			w:      expr.Width64,
			result: valUint(0, expr.Width64),
		},
	}
	testBinaryOp(t, tests, expreval.Lsh)
}

func TestRsh(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "byte_shift",
			v1:     expreval.Value{1, 2, 3, 4},
			v2:     valUint(16, expr.Width8),
			w:      expr.Width32,
			result: expreval.Value{3, 4, 0, 0},
		},
		{
			name:   "bit_shift",
			v1:     expreval.Value{0xff, 0, 0x0f, 0x2},
			v2:     valUint(3, expr.Width8),
			w:      expr.Width32,
			result: expreval.Value{0x1f, 0xe0, 0x41, 0},
		},
		{
			name:   "shift_cut_to_zero",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     expreval.Value{0, 0, 0, 0, 0, 0, 0, 0, 1},
			w:      expr.Width64,
			result: valUint(math.MaxUint64, expr.Width64),
		},
		{
			name:   "shift_above_uint64",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     expreval.Value{0, 0, 0, 0, 0, 0, 0, 0, 1},
			w:      expr.Width128,
			result: valUint(0, expr.Width128),
		},
		{
			name:   "byte_shift_equals_width",
			v1:     valUint(0xf423b129, expr.Width64),
			v2:     valUint(65, expr.Width8),
			w:      expr.Width64,
			result: valUint(0, expr.Width64),
		},
		{
			name:   "byte_shift_above_width",
			v1:     valUint(0xf423b129, expr.Width64),
			v2:     valUint(9*8+2, expr.Width8),
			w:      expr.Width64,
			result: valUint(0, expr.Width64),
		},
	}
	testBinaryOp(t, tests, expreval.Rsh)
}

func TestMul(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "single_byte",
			v1:     expreval.Value{255},
			v2:     expreval.Value{255},
			w:      2,
			result: expreval.Value{0x01, 0xfe},
		},
		{
			name:   "two_bytes",
			v1:     valUint(46576, expr.Width16),
			v2:     valUint(12344, expr.Width16),
			w:      4,
			result: valUint(46576*12344, expr.Width32),
		},
		{
			name:   "zero_extend",
			v1:     expreval.Value{54},
			v2:     valUint(4454466, expr.Width32),
			w:      expr.Width32,
			result: valUint(54*4454466, expr.Width32),
		},
		{
			name:   "cut_upper_bytes",
			v1:     expreval.Value{0, 0, 5, 6},
			v2:     expreval.Value{45, 67},
			w:      expr.Width16,
			result: expreval.Value{0, 0},
		},
		{
			name:   "overflow",
			v1:     expreval.Value{0, 0, 0xff, 0x1},
			v2:     expreval.Value{0x5, 0x7, 0x23, 0x71},
			w:      expr.Width32,
			result: expreval.Value{0x0, 0x0, 0xfb, 0x2},
		},
	}
	testBinaryOp(t, tests, expreval.Mul)
}

func TestDivMod(t *testing.T) {
	tests := []struct {
		name string
		v1   expreval.Value
		v2   expreval.Value
		w    expr.Width
		div  expreval.Value
		mod  expreval.Value
	}{
		{
			name: "single_byte",
			v1:   expreval.Value{183},
			v2:   expreval.Value{23},
			w:    expr.Width8,
			div:  expreval.Value{183 / 23},
			mod:  expreval.Value{183 % 23},
		},
		{
			name: "two_bytes",
			v1:   valUint(54678, expr.Width16),
			v2:   valUint(3345, expr.Width16),
			w:    expr.Width16,
			div:  valUint(54678/3345, expr.Width16),
			mod:  valUint(54678%3345, expr.Width16),
		},
		{
			name: "greater_divisor",
			v1:   expreval.Value{54},
			v2:   valUint(4454466, expr.Width32),
			w:    expr.Width32,
			div:  valUint(0, expr.Width32),
			mod:  valUint(54, expr.Width32),
		},
		{
			name: "cut_upper_bytes",
			v1:   valUint(99340, expr.Width32),
			v2:   valUint(56, expr.Width8),
			w:    expr.Width16,
			div:  valUint((99340&0xffff)/56, expr.Width16),
			mod:  valUint((99340&0xffff)%56, expr.Width16),
		},
		{
			name: "div_same_numbers",
			v1:   valUint(0x5b3c3d7a, expr.Width32),
			v2:   valUint(0x5b3c3d7a, expr.Width32),
			w:    expr.Width32,
			div:  expreval.Value{1, 0, 0, 0},
			mod:  expreval.Value{0, 0, 0, 0},
		},
		{
			name: "div_by_zero",
			v1:   valUint(293445, expr.Width32),
			v2:   expreval.Value{0},
			w:    expr.Width32,
			div:  expreval.Value{0xff, 0xff, 0xff, 0xff},
			mod:  valUint(293445, expr.Width32),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("div", func(t *testing.T) {
				result := expreval.Div(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.div, result)
			})

			t.Run("mod", func(t *testing.T) {
				result := expreval.Mod(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.mod, result)
			})
		})
	}
}

func TestBitwise(t *testing.T) {
	tests := []struct {
		name string
		v1   expreval.Value
		v2   expreval.Value
		w    expr.Width

		and expreval.Value
		or  expreval.Value
		xor expreval.Value
	}{
		{
			name: "multi_byte",
			v1:   expreval.Value{0xfb, 0x45, 0x43, 0x2b},
			v2:   expreval.Value{0xde, 0x04, 0x22, 0x9f},
			w:    expr.Width32,

			and: expreval.Value{0xda, 0x04, 0x02, 0x0b},
			or:  expreval.Value{0xff, 0x45, 0x63, 0xbf},
			xor: expreval.Value{0xfb ^ 0xde, 0x45 ^ 0x04, 0x43 ^ 0x22, 0x2b ^ 0x9f},
		},
		{
			name: "cut_and_extend",
			v1:   expreval.Value{0x55, 0x23, 0xb2, 0x3a, 0x24, 0x6f, 0x34, 0xbb},
			v2:   expreval.Value{0x56, 0xb3},
			w:    expr.Width32,

			and: expreval.Value{0x54, 0x23, 0, 0},
			or:  expreval.Value{0x57, 0xb3, 0xb2, 0x3a},
			xor: expreval.Value{0x55 ^ 0x56, 0x23 ^ 0xb3, 0xb2 ^ 0, 0x3a ^ 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("and", func(t *testing.T) {
				result := expreval.And(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.and, result)
			})
			t.Run("or", func(t *testing.T) {
				result := expreval.Or(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.or, result)
			})
			t.Run("xor", func(t *testing.T) {
				result := expreval.Xor(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.xor, result)
			})
		})
	}
}