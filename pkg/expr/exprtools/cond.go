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
	arg1Max := expr.NewCond(expr.Ltu, arg1, ones, exprFalse, exprTrue, w)
	equal := expr.NewCond(expr.Ltu, arg2, ones, lessThenPlusOne, arg1Max, w)

	return expr.NewCond(expr.Ltu, arg1, arg2, exprFalse, equal, w)
}

// Lts creates an expression which returns exprTrue of width w if arg1 is less
// then to arg2 in signed integer comparison. Otherwise, this expression returns
// exprFalse of width w.
func Lts(arg1, arg2 expr.Expr, exprTrue, exprFalse expr.Expr, w expr.Width) expr.Expr {
	mask := signBitMask(w)
	sign1, sign2 := BitAnd(arg1, mask, w), BitAnd(arg2, mask, w)
	arg1Abs, arg2Abs := abs(arg1, mask), abs(arg2, mask)

	posLess := expr.NewCond(expr.Ltu, arg1Abs, arg2Abs, exprTrue, exprFalse, w)
	negLess := expr.NewCond(expr.Ltu, arg2Abs, arg1Abs, exprTrue, exprFalse, w)
	signEqLess := expr.NewCond(expr.Ltu, expr.Zero, sign1, negLess, posLess, w)

	// If sign1 < sign2, then arg1 is positive and arg2 is negative.
	// Sidenote: In such a case sign1 is zero.
	signNeqLess := expr.NewCond(expr.Ltu, sign1, sign2, exprFalse, exprTrue, w)

	// It's important to keep in mind that sign1 and sign2 have the same
	// width so the sign bit is in the same position.
	signDiff := BitXor(sign1, sign2, w)
	return expr.NewCond(expr.Ltu, expr.Zero, signDiff, signNeqLess, signEqLess, w)

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
	return Lts(arg1, arg2, exprTrue, eq, w)
}
