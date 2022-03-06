package model

import "math"

// Address represents an arbitrary memory address in a program (user-space)
// address space.
//
// We might use plain old uint64 to represent any memory address and we would be
// most likely fine for following 10 years. On the other hand given that RISC-V
// already has 128bit instruction set described, it makes sense to introduce
// this alias and to make the code more variable in the future.
//
// For future compatibility, it's guaranteed that this type will be always an
// unsigned integer. What can change in the future is its bit width, which will
// be most likely expanded to 128 bits.
type Address = uint64

const (
	// MinAddress is the smallest value Address is able to represent.
	MinAddress Address = 0
	// MaxAddress is the biggest value Address is able to represent.
	MaxAddress Address = math.MaxUint64
)
