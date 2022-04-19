package exprtools

import (
	"mltwist/pkg/expr"
)

// NewWidthGadget creates an expression with width w, which wraps e and produces
// the same value as e. This is useful to change expression width.
//
// Even though expr package provides many ways how to do this kind of operation,
// we decided to define the one. Further expr optimization can algorithms might
// be optimized to recognize this particular gadget.
func NewWidthGadget(e expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Add, e, expr.Zero, w)
}

func widthGadget(e expr.Expr) bool {
	b, ok := e.(expr.Binary)
	if !ok || b.Op() != expr.Add {
		return false
	}

	arg2, ok := b.Arg2().(expr.Const)
	if !ok {
		return false
	}

	return arg2.Equal(expr.Zero)
}

// WidthGadgetArg checks if e is width gadget and returns it's argument and true
// in case it is. In case e is not width gadget, the first return value is
// undefined and the later one is false.
func WidthGadgetArg(e expr.Expr) (expr.Expr, bool) {
	if widthGadget(e) {
		b := e.(expr.Binary)
		return b.Arg1(), true

	}

	return nil, false
}
