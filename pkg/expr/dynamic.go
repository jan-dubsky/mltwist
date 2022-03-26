package expr

import "decomp/pkg/expr/internal/value"

var _ Expr = &dynamic{}

type dynamic struct {
	// The structure must be non-empty, as all empty instructions have the
	// same (dummy) address in Go programs. This is undesired as then all
	// dynamic values would equal each other, as they would have the same
	// address and same (empty) content.
	_ uint8
}

func NewDynamic() Expr { return &dynamic{} }

func (d *dynamic) value() (value.Value, bool) { return value.Value{}, false }
