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
	case expr.Nand:
		return expreval.Nand
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
	case expr.Ltu:
		return expreval.Ltu
	case expr.Lts:
		return expreval.Lts
	default:
		panic(fmt.Sprintf("unknown condition type: %v", c))
	}
}

func condEval(c expr.Condition, c1 expr.Const, c2 expr.Const, w expr.Width) bool {
	f := condEvalFunc(c)

	v1, v2 := expreval.ParseConst(c1), expreval.ParseConst(c2)
	return f(v1, v2, w)
}

// ConstFold replaces all expressions of constants by a constant expression
// transitively in a whole expression subtree of ex.
//
// The returned expression doesn't containt any unnecessary width gadgets as
// PurgeWidthGadgets is applied at the end of constant folding process.
//
// In order to perform static analysis of an expression, we have to be able to
// identify static expressions and evaluate them. Constant fold on constant
// naturally produces the same constant. Constant fold of an arbitrary binary
// expression with both constant arguments is replaced by constant corresponding
// to the operation result. Constant fold of register load if the register load
// itself. Constant fold of a memory load is the same memory load, but with it's
// address expression constant-folded. The most complicated expression to
// constant fold is conditional expression, which first constant-folds it's
// arguments for condition evaluation. If both those arguments are constants,
// the result is true of false expression constant-folded depending on the
// condition being true or false. If the condition is not constant-foldable, the
// result of const fold is the same constant expression but with all its
// arguments constant-folded.
//
// This function tries to reuse subtrees of an expression tree rooted in ex as
// much as possible to minimize memory foodprint of the application. In other
// words, if a leaf sub-tree of the expression tree cannot be constant-folded,
// it's reused in the new tree. An edge case of this is that return value can
// equal ex if there are no expressions which can be constant folded.
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
			if condEval(e.Cond(), c1, c2, e.Width()) {
				res, _ = constFold(e.ExprTrue())
			} else {
				res, _ = constFold(e.ExprFalse())
			}

			return SetWidth(res, e.Width()), true
		}

		t, changedTrue := constFold(e.ExprTrue())
		f, changedFalse := constFold(e.ExprFalse())

		if !(changedArg1 || changedArg2 || changedTrue || changedFalse) {
			return ex, false
		}
		return expr.NewCond(e.Cond(), arg1, arg2, t, f, e.Width()), true
	case expr.Const:
		return ex, false
	case expr.MemLoad:
		addr, changedAddr := constFold(e.Addr())

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
