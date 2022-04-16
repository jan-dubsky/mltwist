package expreval

import (
	"decomp/pkg/expr"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type binaryOpTest struct {
	name   string
	v1     value
	v2     value
	w      expr.Width
	result value
}

type binaryOpFunc func(v1 value, v2 value, w expr.Width) value

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
			v1:     value{255},
			v2:     value{1},
			w:      expr.Width32,
			result: value{0, 1, 0, 0},
		},
		{
			name:   "carry_chain",
			v1:     value{255, 255, 255, 255, 255},
			v2:     value{7, 67, 23, 43},
			w:      expr.Width64,
			result: value{6, 67, 23, 43, 0, 1, 0, 0},
		},
		{
			name:   "cut_upper_bytes",
			v1:     value{5, 6, 7, 8},
			v2:     value{34, 255},
			w:      expr.Width16,
			result: value{39, 5},
		},
	}
	testBinaryOp(t, tests, add)
}

func TestSub(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "without_carry",
			v1:     value{45, 67, 89, 8},
			v2:     value{45, 45, 84, 8},
			w:      expr.Width32,
			result: value{0, 22, 5, 0},
		},
		{
			name:   "carry_chain",
			v1:     value{98, 93, 34, 67, 89, 78},
			v2:     value{99, 245, 88, 98, 0, 12},
			w:      expr.Width64,
			result: value{255, 103, 201, 224, 88, 66, 0, 0},
		},
		{
			name:   "cut_upper_bytes",
			v1:     value{5, 6, 7, 8},
			v2:     value{34, 254},
			w:      expr.Width16,
			result: value{227, 7},
		},
	}
	testBinaryOp(t, tests, sub)
}

func TestLsh(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "byte_shift",
			v1:     value{1, 2, 3, 4},
			v2:     valUint(16, expr.Width8),
			w:      expr.Width32,
			result: value{0, 0, 1, 2},
		},
		{
			name:   "bit_shift",
			v1:     value{0xff, 0, 0x0f, 0x2},
			v2:     valUint(3, expr.Width8),
			w:      expr.Width32,
			result: value{0xf8, 7, 0x78, 0x10},
		},
		{
			name:   "shift_cut_to_zero",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     value{0, 0, 0, 0, 0, 0, 0, 0, 1},
			w:      expr.Width64,
			result: valUint(math.MaxUint64, expr.Width64),
		},
		{
			name:   "shift_above_uint64",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     value{0, 0, 0, 0, 0, 0, 0, 0, 1},
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
	testBinaryOp(t, tests, lsh)
}

func TestRsh(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "byte_shift",
			v1:     value{1, 2, 3, 4},
			v2:     valUint(16, expr.Width8),
			w:      expr.Width32,
			result: value{3, 4, 0, 0},
		},
		{
			name:   "bit_shift",
			v1:     value{0xff, 0, 0x0f, 0x2},
			v2:     valUint(3, expr.Width8),
			w:      expr.Width32,
			result: value{0x1f, 0xe0, 0x41, 0},
		},
		{
			name:   "shift_cut_to_zero",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     value{0, 0, 0, 0, 0, 0, 0, 0, 1},
			w:      expr.Width64,
			result: valUint(math.MaxUint64, expr.Width64),
		},
		{
			name:   "shift_above_uint64",
			v1:     valUint(math.MaxUint64, expr.Width64),
			v2:     value{0, 0, 0, 0, 0, 0, 0, 0, 1},
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
	testBinaryOp(t, tests, rsh)
}

func TestMul(t *testing.T) {
	tests := []binaryOpTest{
		{
			name:   "single_byte",
			v1:     value{255},
			v2:     value{255},
			w:      2,
			result: value{0x01, 0xfe},
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
			v1:     value{54},
			v2:     valUint(4454466, expr.Width32),
			w:      expr.Width32,
			result: valUint(54*4454466, expr.Width32),
		},
		{
			name:   "cut_upper_bytes",
			v1:     value{0, 0, 5, 6},
			v2:     value{45, 67},
			w:      expr.Width16,
			result: value{0, 0},
		},
		{
			name:   "overflow",
			v1:     value{0, 0, 0xff, 0x1},
			v2:     value{0x5, 0x7, 0x23, 0x71},
			w:      expr.Width32,
			result: value{0x0, 0x0, 0xfb, 0x2},
		},
	}
	testBinaryOp(t, tests, mul)
}

func TestDivMod(t *testing.T) {
	tests := []struct {
		name string
		v1   value
		v2   value
		w    expr.Width
		div  value
		mod  value
	}{
		{
			name: "single_byte",
			v1:   value{183},
			v2:   value{23},
			w:    expr.Width8,
			div:  value{183 / 23},
			mod:  value{183 % 23},
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
			v1:   value{54},
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
			div:  value{1, 0, 0, 0},
			mod:  value{0, 0, 0, 0},
		},
		{
			name: "div_by_zero",
			v1:   valUint(293445, expr.Width32),
			v2:   value{0},
			w:    expr.Width32,
			div:  value{0xff, 0xff, 0xff, 0xff},
			mod:  valUint(293445, expr.Width32),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("div", func(t *testing.T) {
				result := div(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.div, result)
			})

			t.Run("mod", func(t *testing.T) {
				result := mod(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.mod, result)
			})
		})
	}
}

func TestBitwise(t *testing.T) {
	tests := []struct {
		name string
		v1   value
		v2   value
		w    expr.Width

		and value
		or  value
		xor value
	}{
		{
			name: "multi_byte",
			v1:   value{0xfb, 0x45, 0x43, 0x2b},
			v2:   value{0xde, 0x04, 0x22, 0x9f},
			w:    expr.Width32,

			and: value{0xda, 0x04, 0x02, 0x0b},
			or:  value{0xff, 0x45, 0x63, 0xbf},
			xor: value{0xfb ^ 0xde, 0x45 ^ 0x04, 0x43 ^ 0x22, 0x2b ^ 0x9f},
		},
		{
			name: "cut_and_extend",
			v1:   value{0x55, 0x23, 0xb2, 0x3a, 0x24, 0x6f, 0x34, 0xbb},
			v2:   value{0x56, 0xb3},
			w:    expr.Width32,

			and: value{0x54, 0x23, 0, 0},
			or:  value{0x57, 0xb3, 0xb2, 0x3a},
			xor: value{0x55 ^ 0x56, 0x23 ^ 0xb3, 0xb2 ^ 0, 0x3a ^ 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Run("and", func(t *testing.T) {
				result := and(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.and, result)
			})
			t.Run("or", func(t *testing.T) {
				result := or(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.or, result)
			})
			t.Run("xor", func(t *testing.T) {
				result := xor(tt.v1, tt.v2, tt.w)
				require.Equal(t, tt.xor, result)
			})
		})
	}
}
