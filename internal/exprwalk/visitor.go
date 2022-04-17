package exprwalk

import "decomp/pkg/expr"

type ExprWalker interface {
	Binary(e expr.Binary) error
	Cond(e expr.Cond) error
	Const(e expr.Const) error
	RegLoad(e expr.RegLoad) error
	MemLoad(e expr.MemLoad) error
}

var _ ExprWalker = EmptyExprWalker{}

type EmptyExprWalker struct{}

func (EmptyExprWalker) Binary(e expr.Binary) error   { return nil }
func (EmptyExprWalker) Cond(e expr.Cond) error       { return nil }
func (EmptyExprWalker) Const(e expr.Const) error     { return nil }
func (EmptyExprWalker) RegLoad(e expr.RegLoad) error { return nil }
func (EmptyExprWalker) MemLoad(e expr.MemLoad) error { return nil }

var _ EffectWalker = EmptyEffectWalker{}

type EffectWalker interface {
	MemStore(e expr.MemStore) error
	RegStore(e expr.RegStore) error
}

type EmptyEffectWalker struct{}

func (EmptyEffectWalker) MemStore(e expr.MemStore) error { return nil }
func (EmptyEffectWalker) RegStore(e expr.RegStore) error { return nil }
