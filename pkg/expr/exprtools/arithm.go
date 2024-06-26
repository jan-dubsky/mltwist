package exprtools

import (
	"fmt"
	"mltwist/pkg/expr"
)

// Negate returns negative (multiplied by -1) integer value of an expression e
// of width w.
func Negate(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Add, BitNot(e, w), expr.One, w)
}

// Sub subtracts e1 and e2 and returns the difference as expression of width w.
//
// If w is greater than e1 or e2 width, unsigned extension will be used to
// extend both e1 and e2 to w.
func Sub(e1, e2 expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Add, e1, Negate(e2, w), w)
}

// abs applies bit mask to e and returns absolute value of e in w bytes. Value
// if w is width of mask which also determines width of the expression produced.
func abs(e expr.Expr, mask expr.Expr) expr.Expr {
	w := mask.Width()
	return expr.NewLess(e, mask, e, Negate(e, w), w)
}

// Abs returns an absolute value of e with width w.
func Abs(e expr.Expr, w expr.Width) expr.Expr {
	mask := signBitMask(w)
	return abs(e, mask)
}

// Ones returns expression of width w filled with all ones.
//
// This function can be used to represent signed -1 value of any width w.
func Ones(w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Nand, expr.Zero, expr.Zero, w)
}

// Mod returns width bits unsigned division reminder of the first argument
// divided by the second argument. The produced result is of width w.
//
// Please note that signed module can be implemented using unsigned module
// followed by sign resolution logic.
//
// Module by zero doesn't cause any error, but produces result of width w with
// value of the first argument.
func Mod(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	div := expr.NewBinary(expr.Div, e1, e2, w)
	// If e2 is zero, multiple is also zero -> mod will be e1 - 0 == e1.
	multiple := expr.NewBinary(expr.Mul, div, e2, w)
	return Sub(e1, multiple, w)
}

// negativeSignJoin extracts signs out of 2 signed integer expressions and
// returns sign of their integer produce.
//
// Value of 1 is returned in case only e1 or e2 is negative. Otherwise zero is
// returned. The return value is a boolean expression - i.e. its width is always
// one.
func negativeSignJoin(e1 expr.Expr, e2 expr.Expr) expr.Expr {
	sign1 := Bool(IntNegative(e1, e1.Width()))
	sign2 := Bool(IntNegative(e2, e2.Width()))
	return BitXor(sign1, sign2, expr.Width8)
}

// SignedMul performs a signed multiplication of 2 signed integer numbers of
// width w. The width of result is 2*w.
//
// Both signed and unsigned multiplication of 2*b bit numbers always result in
// equal lowest b bits of the result. Because of this, there is no reason to
// implement signed multiplication on w width only as then it could be altered
// by unsigned multiplication. The only case when signed and unsigned
// multiplication differs is in case when result is twice as wide as its
// arguments.
//
// As w is internally doubled and the multiplication must not overflow, the
// maximal allowed value of w is 127. This value will panic for higher value of
// w. Given that typical platform allow multiplication of at most 8 byte values,
// this limitation is mostly theoretical and it should not be never seen in
// practice.
func SignedMul(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	if w > expr.MaxWidth/2 {
		panic(fmt.Errorf("too big signed multiplication width: %d", w))
	}

	e1Ext := SignExtend(e1, expr.ConstFromUint(e1.Width().Bits()-1), 2*w)
	e2Ext := SignExtend(e2, expr.ConstFromUint(e2.Width().Bits()-1), 2*w)

	// Multiplication of k bit results has k lower bits same no matter
	// whether it's signed or unsigned.
	return expr.NewBinary(expr.Mul, e1Ext, e2Ext, 2*w)
}

func signedOp(
	e1 expr.Expr,
	e2 expr.Expr,
	w expr.Width,
	f func(e1, e2 expr.Expr, w expr.Width) expr.Expr,
) expr.Expr {
	unsigned := f(Abs(e1, e1.Width()), Abs(e2, e2.Width()), w)
	return BoolCond(
		negativeSignJoin(e1, e2),
		Negate(unsigned, w),
		unsigned,
		w,
	)
}

// SignedDiv implement signed division of 2 w wide values. The result is as well
// w wide.
//
// Division by zero will result in an expression of width w filled with ones
// (i.e. signed -1).
//
// An overflow can happen of e1 is maximal negative number representable in w
// and e2 is -1. In such a case the result is the same value as e1.
func SignedDiv(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	signedDiv := signedOp(e1, e2, w, func(e1, e2 expr.Expr, w expr.Width) expr.Expr {
		return expr.NewBinary(expr.Div, e1, e2, w)
	})

	// Special case division by zero.
	return BoolCond(
		e2,
		signedDiv,
		Ones(w),
		w,
	)
}

// SignedMod implement signed modulo operation for 2 w-wide values. The result
// is as well w wide.
//// Modulo by zero will result in an expression of width w with the same value
// as e1 (potentially cropped to w bytes).
//
// An overflow can happen of e1 is maximal negative number representable in w
// and e2 is -1. In such a case the result is zero.
func SignedMod(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	return signedOp(e1, e2, w, Mod)
}

// SignExtend implements sign extension of e, where bit at position signBit is
// understood as sign bit. The resulting expression has width w. The result is
// undefined if signBit is higher than width of resulting expression.
//
// SignBit is an index of bit starting with 0. So the least significant bit has
// index zero.
//
// All bits higher than sign bit are set in both positive and negative cases.
// Consequently, SignExpect(5, 1) is 1 ad bit 2 is set to the value of bit 1
// which is zero.
func SignExtend(e expr.Expr, signBit expr.Expr, w expr.Width) expr.Expr {
	signMask := expr.NewBinary(expr.Lsh, expr.One, signBit, w)

	// This is absolute value without sign.
	valueBitsMask := Sub(signMask, expr.One, w)
	valueSignBits := BitAnd(e, valueBitsMask, w)
	signBitsMask := BitNot(valueBitsMask, w)

	return BoolCond(
		BitAnd(e, signMask, w),
		BitOr(e, signBitsMask, w),
		valueSignBits, // positive number,
		w,
	)
}

// RshA arithmetically shifts e to right shift number of bits producing an
// expression of width w.
//
// Highest bits produced by shift are always filled by zeros for non-negative
// number and always filled with ones for negative number.
func RshA(e expr.Expr, shift expr.Expr, w expr.Width) expr.Expr {
	ones := Ones(w)
	shiftedMask := expr.NewBinary(expr.Rsh, ones, shift, w)
	addedBitMask := Sub(ones, shiftedMask, w)

	rsh := expr.NewBinary(expr.Rsh, e, shift, w)
	rshNeg := BitOr(rsh, addedBitMask, w)

	// Note that negative numbers are greater then positive in unsigned
	// comparison.
	return expr.NewLess(
		e,
		signBitMask(w),
		rsh,
		rshNeg,
		w,
	)
}
