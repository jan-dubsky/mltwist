package exprtools

import "mltwist/pkg/expr"

// Bool converts any nonzero expression e into expression with value 1. If e is
// zero, value of zero is returned. The return value has always width 1.
func Bool(e expr.Expr) expr.Expr {
	return NewWidthGadget(
		// e equal 0 if it's less than one as unsigned integers.
		expr.NewLess(
			e,
			expr.One,
			expr.Zero,
			expr.One,
			e.Width(),
		),
		expr.Width8,
	)
}

// Not implements C-like boolean not operation.
//
// If e is nonzero, this function return zero expression. Otherwise it returns
// expression with value one. Width of returned expression is always 1.
func Not(e expr.Expr) expr.Expr {
	return NewWidthGadget(
		// e equal 0 if it's less than one as unsigned integers.
		expr.NewLess(
			e,
			expr.One,
			expr.One,
			expr.Zero,
			e.Width(),
		),
		expr.Width8,
	)
}

// BoolCond implements C-like condition. Of boolExpr is nonzero trueExpr is
// returned. Otherwise falseExpr is returned. In both true and false cases, the
// result width is w.
func BoolCond(boolExpr expr.Expr, trueExpr, falseExpr expr.Expr, w expr.Width) expr.Expr {
	// In normal logic nonzero is true, but we use boolExpr == 0 condition
	// so we have to swap true and false expression.
	return expr.NewLess(
		expr.Zero,
		boolExpr,
		trueExpr,
		falseExpr,
		w,
	)
}
