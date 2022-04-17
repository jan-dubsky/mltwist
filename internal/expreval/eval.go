package expreval

import (
	"decomp/pkg/expr"
	"fmt"
)

func binaryEvalFunc(op expr.BinaryOp) func(v1 value, v2 value, w expr.Width) value {
	switch op {
	case expr.Add:
		return add
	case expr.Sub:
		return sub
	case expr.Lsh:
		return lsh
	case expr.Rsh:
		return rsh
	case expr.Mul:
		return mul
	case expr.Div:
		return div
	case expr.Mod:
		return mod
	case expr.And:
		return and
	case expr.Or:
		return or
	case expr.Xor:
		return xor
	default:
		panic(fmt.Sprintf("unknown binary operation: %v", op))
	}
}

func binaryEval(op expr.BinaryOp, c1 expr.Const, c2 expr.Const, w expr.Width) expr.Const {
	f := binaryEvalFunc(op)

	v1, v2 := parseValueConst(c1), parseValueConst(c2)
	return f(v1, v2, w).castConst(w)
}

func condEvalFunc(c expr.Condition) func(v1 value, v2 value, w expr.Width) bool {
	switch c {
	case expr.Eq:
		return eq
	case expr.Ltu:
		return ltu
	case expr.Leu:
		return leu
	case expr.Lts:
		return lts
	case expr.Les:
		return les
	default:
		panic(fmt.Sprintf("unknown condition type: %v", c))
	}
}

func condEval(c expr.Condition, c1 expr.Const, c2 expr.Const, w expr.Width) bool {
	f := condEvalFunc(c)

	v1, v2 := parseValueConst(c1), parseValueConst(c2)
	return f(v1, v2, w)
}

func setExprWidth(e expr.Expr, w expr.Width) expr.Expr {
	if e.Width() == w {
		return e
	}

	if c, ok := e.(expr.Const); ok {
		return expr.NewConst(c.Bytes(), w)
	}

	return expr.NewBinary(expr.Add, e, expr.Zero, w)
}

func EvalConst(ex expr.Expr) expr.Expr {
	switch e := ex.(type) {
	case expr.Binary:
		arg1, arg2 := EvalConst(e.Arg1()), EvalConst(e.Arg2())
		c1, ok1 := arg1.(expr.Const)
		c2, ok2 := arg2.(expr.Const)
		if !ok1 || !ok2 {
			return expr.NewBinary(e.Op(), arg1, arg2, e.Width())
		}

		return binaryEval(e.Op(), c1, c2, e.Width())
	case expr.Cond:
		arg1, arg2 := EvalConst(e.Arg1()), EvalConst(e.Arg2())
		c1, ok1 := arg1.(expr.Const)
		c2, ok2 := arg2.(expr.Const)
		if !ok1 || !ok2 {
			t, f := EvalConst(e.ExprTrue()), EvalConst(e.ExprFalse())
			return expr.NewCond(e.Condition(), arg1, arg2, t, f, e.Width())
		}

		var res expr.Expr
		if condEval(e.Condition(), c1, c2, e.Width()) {
			res = EvalConst(e.ExprTrue())
		} else {
			res = EvalConst(e.ExprFalse())
		}

		return setExprWidth(res, e.Width())
	case expr.Const:
		return e
	case expr.MemLoad:
		return e
	case expr.RegLoad:
		return e
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}
