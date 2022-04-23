package model

import (
	"math"
	"mltwist/pkg/expr"
	"unsafe"
)

// Addr represents an arbitrary memory address in a program (user-space) address
// space.
//
// We might use plain old uint64 to represent any memory address and we would be
// most likely fine for following 10 years. On the other hand given that RISC-V
// already has 128bit instruction set described, it makes sense to introduce
// this alias and to make the code more variable in the future.
//
// For future compatibility, it's guaranteed that this type will be always an
// unsigned integer and that uint64 will be always castable to this type.
type Addr uint64

const (
	// MinAddress is the smallest value Address is able to represent.
	MinAddress Addr = 0
	// MaxAddress is the biggest value Address is able to represent.
	MaxAddress Addr = math.MaxUint64
)

// AddrExpr is a helper function creating expr.Const with width of Address type
// and value of a.
func AddrExpr(a Addr) expr.Const {
	return expr.NewConstUint(a, expr.Width(unsafe.Sizeof(a)))
}
