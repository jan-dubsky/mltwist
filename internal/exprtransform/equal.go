package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

// Equal compares if ex1 and ex2 are identical expressions.
//
// Expressions ex1 and ex2 are identical if they have the same width, type and
// if all the properties if their respective expression match are equal.
// Equality checks then recursively apply to all arguments of equal.
func Equal(ex1 expr.Expr, ex2 expr.Expr) bool {
	if ex1.Width() != ex2.Width() {
		return false
	}

	switch e1 := ex1.(type) {
	case expr.Binary:
		e2, ok := ex2.(expr.Binary)
		if !ok {
			return false
		}

		return e1.Op() == e2.Op() &&
			Equal(e1.Arg1(), e2.Arg1()) &&
			Equal(e1.Arg2(), e2.Arg2())
	case expr.Less:
		e2, ok := ex2.(expr.Less)
		if !ok {
			return false
		}

		return Equal(e1.Arg1(), e2.Arg1()) &&
			Equal(e1.Arg2(), e2.Arg2()) &&
			Equal(e1.ExprTrue(), e2.ExprTrue()) &&
			Equal(e1.ExprFalse(), e2.ExprFalse())
	case expr.Const:
		e2, ok := ex2.(expr.Const)
		if !ok {
			return false
		}

		return e1.Equal(e2)
	case expr.MemLoad:
		e2, ok := ex2.(expr.MemLoad)
		if !ok {
			return false
		}

		return e1.Key() == e2.Key() &&
			Equal(e1.Addr(), e2.Addr())
	case expr.RegLoad:
		e2, ok := ex2.(expr.RegLoad)
		if !ok {
			return false
		}

		return e1.Equal(e2)
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex1))
	}
}
