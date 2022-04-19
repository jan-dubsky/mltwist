package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

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
	case expr.Cond:
		e2, ok := ex2.(expr.Cond)
		if !ok {
			return false
		}

		return e1.Condition() == e2.Condition() &&
			Equal(e1.Arg1(), e2.Arg1()) &&
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
