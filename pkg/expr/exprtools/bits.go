package exprtools

import "mltwist/pkg/expr"

// BitNot negates every bit in e and returns result of width w.
func BitNot(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Nand, e, e, w)
}

// BitAnd calculates bitwise AND of e1 and e2 of width w.
func BitAnd(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	nand := expr.NewBinary(expr.Nand, e1, e2, w)
	return BitNot(nand, w)
}

// BitOr calculates bitwise OR of e1 and e2 of width w.
func BitOr(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	not1, not2 := BitNot(e1, w), BitNot(e2, w)
	return expr.NewBinary(expr.Nand, not1, not2, w)
}

// BitXor calculates bitwise XOR (exclusive OR) of e1 and e2 of width w.
func BitXor(e1 expr.Expr, e2 expr.Expr, w expr.Width) expr.Expr {
	// https://electronicsphysics.com/xor-gate-diagram-using-only-nand-or-nor-gate
	nandInputs := expr.NewBinary(expr.Nand, e1, e2, w)
	nand1 := expr.NewBinary(expr.Nand, e1, nandInputs, w)
	nand2 := expr.NewBinary(expr.Nand, e2, nandInputs, w)
	return expr.NewBinary(expr.Nand, nand1, nand2, w)
}
