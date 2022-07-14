package model

import (
	"math"
	"mltwist/pkg/expr"
	"unsafe"
)

// Addr represents an arbitrary memory address in a program (user-space) address
// space.
//
// Memory address can overflow. This implies that the memory model we use is
// "donut" memory, where MaxAddr + 1 == MinAddr.
//
// We might use plain old uint64 to represent any memory address and we would be
// most likely fine for following 10 years. On the other hand given that RISC-V
// already has 128bit instruction set described, it makes sense to introduce
// this alias and to make the code more variable in the future.
//
// For future compatibility, it's guaranteed that this type will be always an
// unsigned integer and that uint64 will be always convertible to this type. On
// the other hand it's not guaranteed that Addr will be always convertible to
// uint64.
type Addr uint64

const (
	// MinAddress is the smallest value Address is able to represent.
	MinAddress Addr = 0
	// MaxAddress is the biggest value Address is able to represent.
	MaxAddress Addr = math.MaxUint64
)

// AddrWidth is expression width capable of capturing any addr value.
const AddrWidth = expr.Width(unsafe.Sizeof(Addr(0)))
