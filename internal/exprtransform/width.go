package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
)

func setWidth(ex expr.Expr, w expr.Width) expr.Expr {
	if ex.Width() == w {
		return ex
	}

	switch e := ex.(type) {
	case expr.Binary, expr.Cond:
		return exprtools.NewWidthGadget(ex, w)
	case expr.Const:
		return expr.NewConst(e.Bytes(), w)
	case expr.MemLoad:
		return exprtools.NewWidthGadget(e, w)
	case expr.RegLoad:
		// Register can contain nonzero bytes above e.Width().
		if w > e.Width() {
			return exprtools.NewWidthGadget(e, w)
		}

		return expr.NewRegLoad(e.Key(), w)
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}

func purgeWidthGadgets(ex expr.Expr) expr.Expr {
	switch e := ex.(type) {
	case expr.Binary:
		e1 := dropUselessWidthGadgets(purgeWidthGadgets(e.Arg1()), e.Width())
		e2 := dropUselessWidthGadgets(purgeWidthGadgets(e.Arg2()), e.Width())

		return expr.NewBinary(e.Op(), e1, e2, e.Width())
	case expr.Cond:
		c1 := dropUselessWidthGadgets(purgeWidthGadgets(e.Arg1()), e.Width())
		c2 := dropUselessWidthGadgets(purgeWidthGadgets(e.Arg2()), e.Width())
		et := dropUselessWidthGadgets(purgeWidthGadgets(e.ExprTrue()), e.Width())
		ef := dropUselessWidthGadgets(purgeWidthGadgets(e.ExprFalse()), e.Width())

		return expr.NewCond(e.Condition(), c1, c2, et, ef, e.Width())
	case expr.Const, expr.MemLoad, expr.RegLoad:
		return e
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}

func dropUselessWidthGadgets(ex expr.Expr, w expr.Width) expr.Expr {
	for {
		g, ok := dropUselessWidthGadget(ex, w)
		if !ok {
			break
		}

		ex = g
	}

	return ex
}

func dropUselessWidthGadget(ex expr.Expr, w expr.Width) (expr.Expr, bool) {
	arg, ok := exprtools.WidthGadgetArg(ex)
	if !ok {
		return nil, false
	}

	// Width gadget shrinks value which can drop some higher bits before
	// following extension to w.
	if w > ex.Width() && ex.Width() < arg.Width() {
		return nil, false
	}

	// Width is growing.
	if w >= ex.Width() && ex.Width() >= arg.Width() {
		return arg, true
	}
	// Width is shrinking.
	if w <= ex.Width() && ex.Width() <= arg.Width() {
		return arg, true
	}
	// Width gadget is biggest of those 3 values -> useless zero extend.
	if w < ex.Width() && ex.Width() > arg.Width() {
		return arg, true
	}

	panic(fmt.Sprintf("unreachable: (%v, %d)", ex, w))
}
