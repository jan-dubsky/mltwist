package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

func newWidthGadget(ex expr.Expr, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Add, ex, expr.Zero, w)
}

func setWidth(ex expr.Expr, w expr.Width) expr.Expr {
	if ex.Width() == w {
		return ex
	}

	switch e := ex.(type) {
	case expr.Binary, expr.Cond:
		return newWidthGadget(ex, w)
	case expr.Const:
		return expr.NewConst(e.Bytes(), w)
	case expr.MemLoad:
		return newWidthGadget(e, w)
	case expr.RegLoad:
		// Register can contain nonzero bytes above e.Width().
		if w > e.Width() {
			return newWidthGadget(e, w)
		}

		return expr.NewRegLoad(e.Key(), w)
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}

func purgeWidthGadgets(ex expr.Expr) expr.Expr {
	switch e := ex.(type) {
	case expr.Binary:
		e1 := dropUselessWidthGadget(purgeWidthGadgets(e.Arg1()), e.Width())
		e2 := dropUselessWidthGadget(purgeWidthGadgets(e.Arg2()), e.Width())

		return expr.NewBinary(e.Op(), e1, e2, e.Width())
	case expr.Cond:
		c1 := dropUselessWidthGadget(purgeWidthGadgets(e.Arg1()), e.Width())
		c2 := dropUselessWidthGadget(purgeWidthGadgets(e.Arg2()), e.Width())
		et := dropUselessWidthGadget(purgeWidthGadgets(e.ExprTrue()), e.Width())
		ef := dropUselessWidthGadget(purgeWidthGadgets(e.ExprFalse()), e.Width())

		return expr.NewCond(e.Condition(), c1, c2, et, ef, e.Width())
	case expr.Const, expr.MemLoad, expr.RegLoad:
		return e
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}

func dropUselessWidthGadget(ex expr.Expr, w expr.Width) expr.Expr {
	for uselessWidthGadget(ex, w) {
		ex = ex.(expr.Binary).Arg1()
	}

	return ex
}

func widthGadget(e expr.Binary) bool {
	return e.Op() == expr.Add && Equal(e.Arg2(), expr.Zero)
}

func uselessWidthGadget(ex expr.Expr, w expr.Width) bool {
	b, ok := ex.(expr.Binary)
	if !ok || !widthGadget(b) {
		return false
	}

	// Width is growing.
	if w >= b.Width() && b.Width() >= b.Arg1().Width() {
		return true
	}
	// Width is shrinking.
	if w <= b.Width() && b.Width() <= b.Arg1().Width() {
		return true
	}
	// Width gadget is biggest of those 3 values -> useless zero extend.
	if w < b.Width() && b.Width() > b.Arg1().Width() {
		return true
	}

	// Width gadget shrinks value which can drop some higher bits before
	// following extension to w.
	return false
}
