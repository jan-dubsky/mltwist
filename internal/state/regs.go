package state

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
)

// RegMap represents a set of registers represented by keys and their respective
// current values.
//
// As register rewrite always changes the whole register value and not only some
// bytes of a register, the whole register value is always represented by a
// single expression. So unlike in Memory where we have to care about
// overlapping writes, we can use a single key-value store to represent any
// state of register file,
type RegMap struct {
	m map[expr.Key]expr.Expr
}

// NewRegMap creates a new instance of RegMap scaled to contain reasonable
// number of registers.
func NewRegMap() *RegMap {
	// We need to set some value as 100 registers (default map size in Go)
	// is too many. So 32 is thumbsucked, but more reasonable value.
	return &RegMap{m: make(map[expr.Key]expr.Expr, 32)}
}

// Load loads register value of register k. The expression returned has width w.
// If k is not present in the register map, this function return (nil, false).
func (m *RegMap) Load(k expr.Key, w expr.Width) (expr.Expr, bool) {
	if e, ok := m.m[k]; ok {
		return exprtransform.SetWidth(e, w), true
	}

	return nil, false
}

// Store writes e of width w to register k.
func (m *RegMap) Store(k expr.Key, e expr.Expr, w expr.Width) {
	m.m[k] = exprtransform.SetWidth(e, w)
}

// Len returns number of registers stored in the map.
func (m *RegMap) Len() int { return len(m.m) }

// Values lists all register keys and their respective values stored in the map.
//
// Warning: The map returned by this function must be treated as read-only.
// Modification of the return value might result in undefined behaviour.
func (m *RegMap) Values() map[expr.Key]expr.Expr { return m.m }
