package exprtools

import (
	"decomp/pkg/expr"
	"fmt"
)

// Negate returns negative (multiplied by -1) integer value of an expression e
// and width w.
func Negate(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Sub, expr.Zero, e, w)
}

// Abs returns an absolute value of e with width w.
func Abs(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewCond(
		expr.Les,
		expr.Zero,
		e,
		e,
		Negate(e, w),
		w,
	)
}

// negativeSignJoin extracts signs out of 2 signed integer expressions and
// returns sign of their integer produce.
//
// Value of 1 is returned in case only e1 or e2 is negative. Otherwise zero is
// returned. The return value is a boolean expression - i.e. its width is always
// one.
func negativeSignJoin(e1 expr.Expr, e2 expr.Expr) expr.Expr {
	w1, w2 := e1.Width(), e2.Width()
	sign1, sign2 := Bool(IntNegative(e1, w1)), Bool(IntNegative(e2, w2))
	return expr.NewBinary(expr.Xor, sign1, sign2, expr.Width8)
}

// SignedMul performs a signed multiplication of 2 signed integer numbers of
// width w. The width of result is 2*w.
//
// Both signed and unsigned multiplication of 2 b bit numbers always result in
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

	mul := expr.NewBinary(expr.Mul, Abs(e1, w), Abs(e2, w), 2*w)
	return BoolCond(
		negativeSignJoin(e1, e2),
		Negate(mul, 2*w),
		mul,
		2*w,
	)
}

func signedOp(op expr.BinaryOp, e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	unsigned := expr.NewBinary(op, Abs(e1, e1.Width()), Abs(e2, e2.Width()), w)
	return BoolCond(
		negativeSignJoin(e1, e2),
		Negate(unsigned, w),
		unsigned,
		w,
	)
}

// SignedDiv implement signed division of 2 w wide values. The result is as well
// w wide.
func SignedDiv(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	return signedOp(expr.Div, e1, e2, w)
}

// SignedDiv implement signed module operation for 2 w wide values. The result
// is as well w wide.
func SignedMod(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	return signedOp(expr.Mod, e1, e2, w)
}

// Ones returns expression of width w filled with all ones.
//
// This function can be used to represent signed -1 value of any width w.
func Ones(w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Sub, expr.Zero, expr.One, 2)
}

// SignExtend implements sign extension of e, where bit at position signBit is
// understood as sign bit. The resulting expression has width w. The result is
// undefined if signBit is higher than width of resulting expression.
func SignExtend(e expr.Expr, signBit expr.Expr, w expr.Width) expr.Expr {
	signMask := expr.NewBinary(expr.Lsh, expr.One, signBit, w)

	validBitMask := expr.NewBinary(expr.Sub, signMask, expr.One, w)
	validBits := expr.NewBinary(expr.And, e, validBitMask, w)

	signBitMask := expr.NewBinary(expr.Xor, Ones(w), validBitMask, w)
	negative := expr.NewBinary(expr.Or, validBits, signBitMask, w)

	return BoolCond(
		expr.NewBinary(expr.And, e, signMask, w),
		negative,
		validBits, // positive number,
		w,
	)
}

func RshA(e expr.Expr, shift expr.Expr, w expr.Width) expr.Expr {
	rsh := expr.NewBinary(expr.Rsh, e, shift, w)

	signBitOrigPos := expr.NewConstUint(w.Bits()-1, w) // last bit in W bit number.
	signBitShiftedPos := expr.NewBinary(expr.Sub, signBitOrigPos, shift, w)
	sextRsh := SignExtend(rsh, signBitShiftedPos, w)

	negativeRsh := expr.NewCond(
		expr.Leu,
		shift,
		expr.NewConstUint(w.Bits(), w),
		sextRsh,
		Ones(w),
		w,
	)

	return expr.NewCond(
		expr.Lts,
		e,
		expr.Zero,
		negativeRsh,
		rsh,
		w,
	)
}
