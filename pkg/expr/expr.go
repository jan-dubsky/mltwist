package expr

import "decomp/pkg/expr/internal/value"

type Expr interface {
	value() (value.Value, bool)
}

func equalValue(e1 Expr, e2 Expr) bool {
	v1, ok := e1.value()
	if !ok {
		return false
	}

	v2, ok := e2.value()
	if !ok {
		return false
	}

	return value.Equal(v1, v2)
}

func Equal(e1 Expr, e2 Expr) bool {
	if e1 == e2 {
		return true
	}

	if ok := equalValue(e1, e2); ok {
		return true
	}

	// TODO: Symbolic equality.
	return false
}
