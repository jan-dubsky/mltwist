package exprtools

import "decomp/pkg/expr"

func Bool(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewCond(
		expr.Eq,
		e,
		expr.Zero,
		expr.Zero,
		expr.One,
		w,
	)
}

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
