package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
)

// SetWidth changes width of ex to w.
//
// If expression width can be changed (for example if ex is expr.Const), a new
// expression of the same type as ex is created with width w. If width cannot be
// changed trivially (for example ex is expr.Binary), then
// exprtools.NewWidthGadget is used to change the width.
//
// Given that this function lacks information about context (i.e. consumer) of
// ex, it's not able to perform any context-based optimization. Consequently the
// expression produced by any algorithm using this method might contain useless
// width gadgets. For this reason, it's highly recommended to use
// PurgeWidthGadgets method on the resultion expression.
func SetWidth(ex expr.Expr, w expr.Width) expr.Expr {
	if e, ok := setWidth(ex, w); ok {
		return e
	}

	return exprtools.NewWidthGadget(ex, w)
}

// SetWidthConst changes width of ex to w and returns the new constant. If ex is
// already a constant of width w, ex is returned.
func SetWidthConst(ex expr.Const, w expr.Width) expr.Const {
	if ex.Width() == w {
		return ex
	}
	return expr.NewConst(ex.Bytes(), w)
}

// setWidth tries to directly change ex width to w - i.e. it creates a new
// expression  of the same type but with width w. This function returns (nil,
// false) in case it's not possible to create such a new expression which would
// be equal to ex under any circumstances.
func setWidth(ex expr.Expr, w expr.Width) (expr.Expr, bool) {
	if ex.Width() == w {
		return ex, true
	}

	switch e := ex.(type) {
	case expr.Binary, expr.Cond:
		// We ignore any smart optimization here as it's significantly
		// simpler to now enter gadgets and drop them later in
		// purgeWidthGadgets function.
		return nil, false
	case expr.Const:
		return expr.NewConst(e.Bytes(), w), true
	case expr.MemLoad:
		// By making the load wider/narrower, we might break dependency
		// analysis. Let's show an example:
		//      STORE16 r15 -> r1(2)
		//      LOAD32 r1(0) -> r4
		//      ADD16 r4, r5 -> r8
		// Where <n> in <ins><n> notation represents number of bits
		// operated and <c> in r<x>(<c>) represents offset in bytes.
		//
		// In setup above, the LOAD32 clearly depends on STORE16 as it
		// reads the same memory. But if we made load narrower because
		// we then use only low 16 bits in ADD16, we'd make those 2
		// instructions independent. Consequently, we are not allowed to
		// shrink memory load width.
		//
		// It's very obvious why we cannot make load wider.
		return nil, false
	case expr.RegLoad:
		// Register can contain nonzero bytes above e.Width().
		if w > e.Width() {
			return nil, false
		}

		// Argument used above for expr.MemLoad doesn't apply here as we
		// don't track dependencies by register bytes as we do with
		// memory, but we use register key to track dependencies.
		return expr.NewRegLoad(e.Key(), w), true
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}

// PurgeWidthGadgets removes all unnecessary width gadgets in an expression
// tree. No no-width gadget expression are not modified in any way.
//
// In many algorithms or transformations, the algorithm doesn't know the full
// context (i.e. consumer) of an expression. This absence of knowledge can
// result in extensive usage of width gadget (defined by exprtools package).
// This function performs context-based analysis of an expression which allows
// it to remove unnecessary width gadgets.
func PurgeWidthGadgets(ex expr.Expr) expr.Expr {
	e, _ := purgeWidthGadgetsKeepWidth(ex)
	return e
}

func purgeWidthGadgetsKeepWidth(ex expr.Expr) (expr.Expr, bool) {
	e, nestedPurged := purgeWidthGadgets(ex)

	var changed bool
	for {
		arg, ok := exprtools.WidthGadgetArg(e)
		if !ok || arg.Width() != e.Width() {
			break
		}

		e = arg
		changed = true
	}

	return e, changed || nestedPurged
}

func purgeWidthGadgets(ex expr.Expr) (expr.Expr, bool) {
	switch e := ex.(type) {
	case expr.Binary:
		e1, changedArg1 := purgeWidthGadgets(e.Arg1())
		e2, changedArg2 := purgeWidthGadgets(e.Arg2())
		e1, prunedArg1 := pruneUselessWidthGadgets(e1, e.Width())
		e2, prunedArg2 := pruneUselessWidthGadgets(e2, e.Width())

		// Performance (allocation) optimization.
		if !(changedArg1 || changedArg2 || prunedArg1 || prunedArg2) {
			return ex, false
		}
		return expr.NewBinary(e.Op(), e1, e2, e.Width()), true
	case expr.Cond:
		c1, changedArg1 := purgeWidthGadgets(e.Arg1())
		c2, changedArg2 := purgeWidthGadgets(e.Arg2())
		et, changedTrue := purgeWidthGadgets(e.ExprTrue())
		ef, changedFalse := purgeWidthGadgets(e.ExprFalse())
		changed := changedArg1 || changedArg2 || changedTrue || changedFalse

		c1, prunedArg1 := pruneUselessWidthGadgets(c1, e.Width())
		c2, prunedArg2 := pruneUselessWidthGadgets(c2, e.Width())
		et, prunedTrue := pruneUselessWidthGadgets(et, e.Width())
		ef, prunedFalse := pruneUselessWidthGadgets(ef, e.Width())

		// Performance (allocation) optimization.
		if !(changed || prunedArg1 || prunedArg2 || prunedTrue || prunedFalse) {
			return ex, false
		}
		return expr.NewCond(e.Condition(), c1, c2, et, ef, e.Width()), true
	case expr.MemLoad:
		// Address keeps its width.
		addr, changedAddr := purgeWidthGadgetsKeepWidth(e.Addr())
		addr, prunedAddr := pruneUselessWidthGadgets(addr, e.Width())

		// Performance (allocation) optimization.
		if !(changedAddr || prunedAddr) {
			return ex, false
		}
		return expr.NewMemLoad(e.Key(), addr, e.Width()), true
	case expr.Const, expr.RegLoad:
		return e, false
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}

func pruneUselessWidthGadgets(ex expr.Expr, w expr.Width) (expr.Expr, bool) {
	var dropped bool
	for {
		g, ok := dropUselessWidthGadget(ex, w)
		if !ok {
			break
		}

		ex = g
		dropped = true
	}

	return ex, dropped
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
