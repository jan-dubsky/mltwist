package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

// ExprReplaceFunc is a function replacing specific expression type.
//
// This function returns a new expression which replaces ex and a boolean
// indicator whether the value returned should be used. If false is returned as
// second argument, the returned expression is ignored. On the other hand if
// second return value is true, the expression returned must be non-nil and the
// returned expression replaces ex.
type ExprReplaceFunc[T expr.Expr] func(ex T) (expr.Expr, bool)

// ReplaceAll replaces every expression e of type T in an expression subtree of
// ex by f(e) and returned a new expression.
//
// This function tries to reuse subtrees of an expression tree rooted in ex as
// much as possible to minimize memory foodprint of the application. In other
// words, if a leaf sub-tree of the expression tree is not changed by f, it's
// reused in the new tree. An edge case of this is that return value can equal
// ex if there were no expression changed by f..
func ReplaceAll[T expr.Expr](ex expr.Expr, f ExprReplaceFunc[T]) expr.Expr {
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

		if changedArg1 || changedArg2 {
			ex, changed = expr.NewBinary(e.Op(), arg1, arg2, e.Width()), true
		}
	case expr.Less:
		arg1, changedArg1 := replaceAll(e.Arg1(), f)
		arg2, changedArg2 := replaceAll(e.Arg2(), f)
		et, changedTrue := replaceAll(e.ExprTrue(), f)
		ef, changedFalse := replaceAll(e.ExprFalse(), f)

		changed = changedArg1 || changedArg2 || changedTrue || changedFalse
		if changed {
			ex = expr.NewLess(arg1, arg2, et, ef, e.Width())
		}

	case expr.MemLoad:
		addr, changedAddr := replaceAll(e.Addr(), f)

		if changedAddr {
			ex, changed = expr.NewMemLoad(e.Key(), addr, e.Width()), true
		}
	case expr.Const, expr.RegLoad:
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}

	e, ok := ex.(T)
	if !ok {
		return ex, changed
	}

	replaced, ok := f(e)
	if !ok {
		return ex, changed
	}
	if replaced == nil {
		panic("function returned nil as new value of an expression")
	}

	return replaced, true
}
