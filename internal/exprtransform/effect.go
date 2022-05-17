package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

// Exprs returns list of expression the effect depends on.
func Exprs(effect expr.Effect) []expr.Expr {
	switch e := effect.(type) {
	case expr.MemStore:
		return []expr.Expr{e.Addr(), e.Value()}
	case expr.RegStore:
		return []expr.Expr{e.Value()}
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", effect))
	}
}

// ExprsMany applies Exprs to every element of effects array and returns
// concatenation of return values into a single array.
func ExprsMany(effects []expr.Effect) []expr.Expr {
	exprs := make([]expr.Expr, 0, len(effects))
	for _, effect := range effects {
		exprs = append(exprs, Exprs(effect)...)
	}

	return exprs
}

// ExprTransformFunc is a function transforming one expression to another.
type ExprTransformFunc func(ex expr.Expr) expr.Expr

// EffectApply applies f to every expr.Expr in effect and produces a new effect
// with the same type and width, but refering transformed expr.Exprs.
func EffectApply(effect expr.Effect, f ExprTransformFunc) expr.Effect {
	switch e := effect.(type) {
	case expr.MemStore:
		return expr.NewMemStore(f(e.Value()), e.Key(), f(e.Addr()), e.Width())
	case expr.RegStore:
		return expr.NewRegStore(f(e.Value()), e.Key(), e.Width())
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", effect))
	}
}

// EffectsApply calls EffectApply for every effect in effects and returns array
// of modified effects. Both order and number of effects is preserved.
func EffectsApply(effects []expr.Effect, f ExprTransformFunc) []expr.Effect {
	applied := make([]expr.Effect, len(effects))
	for i, e := range effects {
		applied[i] = EffectApply(e, f)
	}

	return applied
}
