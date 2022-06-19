package exprtools

import "mltwist/pkg/expr"

// Eq creates an expression which returns exprTrue of width w if arg1 is equal
// to arg2 and returns exprFalse (of width w) otherwise.
func Eq(arg1, arg2 expr.Expr, exprTrue, exprFalse expr.Expr, w expr.Width) expr.Expr {
	ones := Ones(w)

	// Arguments equal implies that arg1 is less then arg2 + 1.
	arg2PlusOne := expr.NewBinary(expr.Add, arg2, expr.One, w)
	lessThenPlusOne := expr.NewCond(expr.Ltu, arg1, arg2PlusOne, exprTrue, exprFalse, w)

	// We have to cover edge case when arg2 is maximal possible value as in
	// such a case the arg2+1 expression overflows.
	arg1NotMax := expr.NewCond(expr.Ltu, arg1, ones, exprFalse, exprTrue, w)
	equal := expr.NewCond(expr.Ltu, arg2, ones, lessThenPlusOne, arg1NotMax, w)

	return expr.NewCond(expr.Ltu, arg1, arg2, exprFalse, equal, w)
}

// Leu creates an expression which evaluates to exprTrue if arg1 is less then
// equal to arg2 in unsigned integer comparison. The width of the comparison and
// the result is w.
func Leu(arg1, arg2 expr.Expr, exprTrue, exprFalse expr.Expr, w expr.Width) expr.Expr {
	eq := Eq(arg1, arg2, exprTrue, exprFalse, w)
	return expr.NewCond(expr.Ltu, arg1, arg2, exprTrue, eq, w)
}

// Les creates an expression which evaluates to exprTrue if arg1 is less then
// equal to arg2 in signed integer comparison. The width of the comparison and
// the result is w.
func Les(arg1, arg2 expr.Expr, exprTrue, exprFalse expr.Expr, w expr.Width) expr.Expr {
	eq := Eq(arg1, arg2, exprTrue, exprFalse, w)
	return expr.NewCond(expr.Lts, arg1, arg2, exprTrue, eq, w)
}
