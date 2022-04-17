package exprwalk

import "decomp/pkg/expr"

type constFinder struct {
	EmptyExprWalker
	nonConst bool
}

func (f *constFinder) RegLoad(e expr.RegLoad) error {
	f.nonConst = true
	return ErrStopWalk
}

func (f *constFinder) MemLoad(e expr.MemLoad) error {
	f.nonConst = true
	return ErrStopWalk
}

func Const(ex expr.Expr) bool {
	f := &constFinder{}
	// The only error returned can be ErrStopWalk
	_ = Expr(f, ex)
	return !f.nonConst
}
