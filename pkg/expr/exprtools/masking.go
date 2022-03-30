package exprtools

import (
	"decomp/pkg/expr"
)

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

func MaskBits(e expr.Expr, bits uint16, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.And, e, bitMask(bits, w), w)
}

func signBitMask(w expr.Width) expr.Expr {
	signBit := w.Bits() - 1
	if signBit < 64 {
		return expr.NewConstUint(uint64(1)<<uint64(signBit), w)
	}

	shift := expr.NewConstUint(signBit, expr.Width16)
	return expr.NewBinary(expr.Lsh, expr.One, shift, w)
}

func IntNegative(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.And, e, signBitMask(w), w)
}

func Crop(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Add, e, expr.Zero, w)
}
