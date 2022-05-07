package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

// FindAll finds all expression of type T in an expression tree of ex.
func FindAll[T expr.Expr](ex expr.Expr) []T { return findAll[T](ex, nil) }

func findAll[T expr.Expr](ex expr.Expr, found []T) []T {
	if e, ok := ex.(T); ok {
		found = append(found, e)
	}

	switch e := ex.(type) {
	case expr.Binary:
		found = findAll(e.Arg1(), found)
		found = findAll(e.Arg2(), found)
	case expr.Cond:
		found = findAll(e.Arg1(), found)
		found = findAll(e.Arg2(), found)
		found = findAll(e.ExprTrue(), found)
		found = findAll(e.ExprFalse(), found)
	case expr.Const, expr.RegLoad:
	case expr.MemLoad:
		found = findAll(e.Addr(), found)
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}

	return found
}
