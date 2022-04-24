package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

type inputExpr interface {
	expr.Const | expr.MemLoad | expr.RegLoad
}

func ReplaceAll[T expr.Expr](ex expr.Expr, f func(curr T) (expr.Expr, bool)) expr.Expr {
	e, _ := replaceAll(ex, f)
	return e
}

func replaceAll[T expr.Expr](
	ex expr.Expr,
	f func(curr T) (expr.Expr, bool),
) (expr.Expr, bool) {
	var changed bool

	switch e := ex.(type) {
	case expr.Binary:
		arg1, changedArg1 := replaceAll(e.Arg1(), f)
		arg2, changedArg2 := replaceAll(e.Arg2(), f)

		// Memory optimization.
		if changedArg1 || changedArg2 {
			ex, changed = expr.NewBinary(e.Op(), arg1, arg2, e.Width()), true
		}
	case expr.Cond:
		arg1, changedArg1 := replaceAll(e.Arg1(), f)
		arg2, changedArg2 := replaceAll(e.Arg2(), f)
		et, changedTrue := replaceAll(e.ExprTrue(), f)
		ef, changedFalse := replaceAll(e.ExprFalse(), f)

		// Memory optimization.
		if changedArg1 || changedArg2 || changedTrue || changedFalse {
			ex = expr.NewCond(e.Condition(), arg1, arg2, et, ef, e.Width())
			changed = true
		}

	case expr.MemLoad:
		addr, changedAddr := replaceAll(e.Addr(), f)

		// Memory optimization.
		if changedAddr {
			ex, changed = expr.NewMemLoad(e.Key(), addr, e.Width()), true
		}
	case expr.Const, expr.RegLoad:
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}

	var ch bool
	if e, ok := ex.(T); ok {
		ex, ch = f(e)
		if ex == nil {
			panic("function returned nil as new value of an expression")
		}
	}

	return ex, changed || ch
}
