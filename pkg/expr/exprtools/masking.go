package exprtools

import (
	"decomp/pkg/expr"
)

// bitMask returns an expression of wirth w which mas ones in positions
// [0..bits] and zeros in all positions above.
func bitMask(bits uint16, w expr.Width) expr.Expr {
	// Optimization for those masks we are able to calculate using
	// in-language features to make the description simpler and possible
	// further execution faster.
	if bits <= 64 {
		return expr.NewConstUint(uint64((1<<bits)-1), w)
	}

	shift := expr.NewConstUint(bits, expr.Width16)
	topBit := expr.NewBinary(expr.Lsh, expr.One, shift, w)
	return expr.NewBinary(expr.Sub, topBit, expr.One, w)
}

// MaskBits returns an expresion of width w with cnt lower bits of e. All higher
// bits are unset.
func MaskBits(e expr.Expr, cnt uint16, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.And, e, bitMask(cnt, w), w)
}

// signBitMask returns a mask with width w with bit at position (w*8)-1 set. In
// other words, the mask returned is a mask of highest bit in the expression.
func signBitMask(w expr.Width) expr.Expr {
	signBit := w.Bits() - 1
	if signBit < 64 {
		return expr.NewConstUint(uint64(1)<<uint64(signBit), w)
	}

	shift := expr.NewConstUint(signBit, expr.Width16)
	return expr.NewBinary(expr.Lsh, expr.One, shift, w)
}

// IntNegative returns a nonzero expression if value of e cropped to w bits is
// negative signed integer. On the other hand if e is a positive integer with w
// bits, the value returned is zero.
//
// To make the expression tree as simple as possible, we don't guarantee the
// expression to be one in case of e being a negative signed integer, but we
// only guarantee to be nonzero. This allows the implementation to use a single
// AND operation in a sign bit. In case a 0/1 result is required, this
// expression can be trivially chained with Bool function.
func IntNegative(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.And, e, signBitMask(w), w)
}

// Crop changes width of e to w.
func Crop(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Add, e, expr.Zero, w)
}
