package exprtools

import (
	"mltwist/pkg/expr"
)

// bitCnt represents a number of bits in an expression. Given that our Width in
// bytes always fits uint8, uint16 is big enough to represent any number of bits
// in an expression.
type BitCnt uint16

// bitMask returns an expression of width w which has ones in positions
// [0..bits) and zeros in position bits and all bits above.
func bitMask(bits BitCnt, w expr.Width) expr.Expr {
	// Optimization for those masks we are able to calculate using
	// in-language features to make the description simpler and possible
	// further execution faster.
	if bits <= 64 {
		return expr.NewConstUint(uint64((1<<bits)-1), w)
	}

	shift := expr.ConstFromUint(bits)
	topBit := expr.NewBinary(expr.Lsh, expr.One, shift, w)
	return Sub(topBit, expr.One, w)
}

// MaskBits returns an expression of width w with cnt lower bits of e. All higher
// bits are unset.
func MaskBits(e expr.Expr, cnt BitCnt, w expr.Width) expr.Expr {
	return BitAnd(e, bitMask(cnt, w), w)
}

// signBitMask returns a mask with width w with bit at position (w*8)-1 set. In
// other words, the mask returned is a mask of highest bit in the expression.
func signBitMask(w expr.Width) expr.Expr {
	signBit := w.Bits() - 1
	if signBit < 64 {
		return expr.NewConstUint(uint64(1)<<uint64(signBit), w)
	}

	shift := expr.ConstFromUint(signBit)
	return expr.NewBinary(expr.Lsh, expr.One, shift, w)
}

// IntNegative returns a nonzero expression if value of e cropped to w bits is
// negative signed integer. On the other hand if e is a positive integer with w
// bits, the value returned is zero.
//
// To make the expression tree as simple as possible, we don't guarantee the
// returned value to be one in case of e being a negative signed integer, but we
// only guarantee it to be nonzero. This allows the implementation to use a
// single AND operation in a sign bit. In case a 0/1 result is required, this
// helper can be trivially chained with Bool function.
func IntNegative(e expr.Expr, w expr.Width) expr.Expr {
	return BitAnd(e, signBitMask(w), w)
}
