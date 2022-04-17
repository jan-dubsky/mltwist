package exprtools

import "mltwist/pkg/expr"

// Bool converts any nonzero expression e into expression with value 1. If e is
// zero, value of zero is returned. The return value has always width 1.
func Bool(e expr.Expr) expr.Expr {
	return expr.NewCond(
		expr.Eq,
		e,
		expr.Zero,
		expr.Zero,
		expr.One,
		expr.Width8,
	)
}

// Not implements C-like boolean not operation.
//
// If e is nonzero, this function return zero expression with width w. Otherwise
// it returns expression with value one and width w.
func Not(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewCond(
		expr.Eq,
		e,
		expr.Zero,
		expr.One,
		expr.Zero,
		w,
	)
}

// BoolCond implements C-like condition. Of boolExpr is nonzero trueExpr is
// returned. Otherwise falseExpr is returned. In both true and false cases, the
// result width is w.
func BoolCond(boolExpr expr.Expr, trueExpr, falseExpr expr.Expr, w expr.Width) expr.Expr {
	// In normal logic nonzero is true, but we use boolExpr == 0 condition
	// so we have to swap true and false expression.
	return expr.NewCond(
		expr.Eq,
		boolExpr,
		expr.Zero,
		falseExpr,
		trueExpr,
		w,
	)
}
