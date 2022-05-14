package state

import (
	"mltwist/pkg/expr"
)

// RegMap represents a set of registers represented by keys and their respective
// current values.
//
// As register rewrite always changes the whole register value and not only some
// bytes of a register, the whole register value if always represented by a
// single expression. So unlike in Memory where we have to care about
// overlapping writes, we can use a single key-value store to represent any
// state of register file,
type RegMap map[expr.Key]expr.Expr

// NewRegMap creates a new instance of RegMap scaled to contain reasonable
// number of registers.
func NewRegMap() RegMap {
	// We need to set some value as 100 registers (default map size in Go)
	// is too many. So 32 is thumbsucked, but more reasonable value.
	return make(RegMap, 32)
}
