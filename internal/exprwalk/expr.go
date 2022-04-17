package exprwalk

import (
	"decomp/pkg/expr"
	"errors"
	"fmt"
)

var ErrStopWalk = fmt.Errorf("exprwalk: stop walk")

func Expr(w ExprWalker, ex expr.Expr) error {
	err := walkExpr(w, ex)
	if errors.Is(err, ErrStopWalk) {
		return nil
	}
	return err
}

func exprs(w ExprWalker, exs ...expr.Expr) error {
	for _, ex := range exs {
		err := walkExpr(w, ex)
		if err != nil {
			return err
		}
	}

	return nil
}

func errExprs(err error, w ExprWalker, exs ...expr.Expr) error {
	if err != nil {
		return err
	}

	return exprs(w, exs...)
}

func walkExpr(w ExprWalker, ex expr.Expr) error {
	switch e := ex.(type) {
	case expr.Binary:
		return errExprs(w.Binary(e), w, e.Arg1(), e.Arg2())
	case expr.Cond:
		return errExprs(w.Cond(e), w, e.Arg1(), e.Arg2(), e.ExprTrue(), e.ExprFalse())
	case expr.Const:
		return w.Const(e)
	case expr.MemLoad:
		return errExprs(w.MemLoad(e), w, e.Addr())
	case expr.RegLoad:
		return w.RegLoad(e)
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %t", ex))
	}
}
