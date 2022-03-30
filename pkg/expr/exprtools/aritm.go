package exprtools

import (
	"decomp/pkg/expr"
	"fmt"
)

func Negate(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Sub, expr.Zero, e, w)
}

func Abs(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewCond(
		expr.Le,
		expr.Zero,
		e,
		e,
		Negate(e, w),
		w,
	)
}

func negativeSignJoin(e1 expr.Expr, e2 expr.Expr) expr.Expr {
	w1, w2 := e1.Width(), e2.Width()
	sign1, sign2 := Bool(IntNegative(e1, w1), w1), Bool(IntNegative(e2, w2), w2)
	return expr.NewBinary(expr.Xor, sign1, sign2, expr.Width8)
}

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

func SignedDiv(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	return signedOp(expr.Div, e1, e2, w)
}

func SignedMod(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	return signedOp(expr.Mod, e1, e2, w)
}
