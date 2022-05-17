package exprtools

import "mltwist/pkg/expr"

// Leu creates an expression which evaluates to exprTrue if arg1 is less then
// equal to arg2 in unsigned integer comparison. The width of the comparison and
// the result is w.
func Leu(arg1, arg2 expr.Expr, exprTrue, exprFalse expr.Expr, w expr.Width) expr.Expr {
	eq := expr.NewCond(expr.Eq, arg1, arg2, exprTrue, exprFalse, w)
	return expr.NewCond(expr.Ltu, arg1, arg2, exprTrue, eq, w)
}

// Les creates an expression which evaluates to exprTrue if arg1 is less then
// equal to arg2 in signed integer comparison. The width of the comparison and
// the result is w.
func Les(arg1, arg2 expr.Expr, exprTrue, exprFalse expr.Expr, w expr.Width) expr.Expr {
	eq := expr.NewCond(expr.Eq, arg1, arg2, exprTrue, exprFalse, w)
	return expr.NewCond(expr.Lts, arg1, arg2, exprTrue, eq, w)
}
