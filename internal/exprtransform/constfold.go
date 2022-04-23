package exprtransform

import (
	"fmt"
	"mltwist/internal/exprtransform/internal/expreval"
	"mltwist/pkg/expr"
)

type binaryFunc func(v1 expreval.Value, v2 expreval.Value, w expr.Width) expreval.Value
type condFunc func(v1 expreval.Value, v2 expreval.Value, w expr.Width) bool

func binaryEvalFunc(op expr.BinaryOp) binaryFunc {
	switch op {
	case expr.Add:
		return expreval.Add
	case expr.Sub:
		return expreval.Sub
	case expr.Lsh:
		return expreval.Lsh
	case expr.Rsh:
		return expreval.Rsh
	case expr.Mul:
		return expreval.Mul
	case expr.Div:
		return expreval.Div
	case expr.Mod:
		return expreval.Mod
	case expr.And:
		return expreval.And
	case expr.Or:
		return expreval.Or
	case expr.Xor:
		return expreval.Xor
	default:
		panic(fmt.Sprintf("unknown binary operation: %v", op))
	}
}

func binaryEval(op expr.BinaryOp, c1 expr.Const, c2 expr.Const, w expr.Width) expr.Const {
	f := binaryEvalFunc(op)

	v1, v2 := expreval.ParseConst(c1), expreval.ParseConst(c2)
	return f(v1, v2, w).Const(w)
}

func condEvalFunc(c expr.Condition) condFunc {
	switch c {
	case expr.Eq:
		return expreval.Eq
	case expr.Ltu:
		return expreval.Ltu
	case expr.Leu:
		return expreval.Leu
	case expr.Lts:
		return expreval.Lts
	case expr.Les:
		return expreval.Les
	default:
		panic(fmt.Sprintf("unknown condition type: %v", c))
	}
}

func condEval(c expr.Condition, c1 expr.Const, c2 expr.Const, w expr.Width) bool {
	f := condEvalFunc(c)

	v1, v2 := expreval.ParseConst(c1), expreval.ParseConst(c2)
	return f(v1, v2, w)
}

func ConstFold(ex expr.Expr) expr.Expr {
	e, _ := constFold(ex)
	return PurgeWidthGadgets(e)
}

func constFold(ex expr.Expr) (expr.Expr, bool) {
	switch e := ex.(type) {
	case expr.Binary:
		arg1, changedArg1 := constFold(e.Arg1())
		arg2, changedArg2 := constFold(e.Arg2())

		c1, ok1 := arg1.(expr.Const)
		c2, ok2 := arg2.(expr.Const)
		if ok1 && ok2 {
			return binaryEval(e.Op(), c1, c2, e.Width()), true
		}

		// Performance (allocation) optimization.
		if !(changedArg1 || changedArg2) {
			return ex, false
		}
		return expr.NewBinary(e.Op(), arg1, arg2, e.Width()), true
	case expr.Cond:
		arg1, changedArg1 := constFold(e.Arg1())
		arg2, changedArg2 := constFold(e.Arg2())

		c1, ok1 := arg1.(expr.Const)
		c2, ok2 := arg2.(expr.Const)
		if ok1 && ok2 {
			var res expr.Expr
			if condEval(e.Condition(), c1, c2, e.Width()) {
				res, _ = constFold(e.ExprTrue())
			} else {
				res, _ = constFold(e.ExprFalse())
			}

			return SetWidth(res, e.Width()), true
		}

		t, changedTrue := constFold(e.ExprTrue())
		f, changedFalse := constFold(e.ExprFalse())

		// Performance (allocation) optimization.
		if !(changedArg1 || changedArg2 || changedTrue || changedFalse) {
			return ex, false
		}
		return expr.NewCond(e.Condition(), arg1, arg2, t, f, e.Width()), true
	case expr.Const:
		return ex, false
	case expr.MemLoad:
		addr, changedAddr := constFold(e.Addr())

		// Performance (allocation) optimization.
		if !changedAddr {
			return ex, false
		}
		return expr.NewMemLoad(e.Key(), addr, e.Width()), true
	case expr.RegLoad:
		return ex, false
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}
