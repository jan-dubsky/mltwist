package expr

import "decomp/pkg/expr/internal/value"

var _ Expr = Negate{}

type Negate struct {
	e Expr
}

func NewNegate(e Expr) Negate {
	return Negate{e: e}
}

func (n Negate) value() (value.Value, bool) {
	v, ok := n.e.value()
	if !ok {
		return value.Value{}, false
	}

	if v.Zero() {
		return v, true
	}

	v = v.Clone()
	v.Negate()
	return v, true
}
