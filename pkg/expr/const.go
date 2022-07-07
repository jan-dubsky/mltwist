package expr

import (
	"bytes"
	"fmt"
	"unsafe"

	"golang.org/x/exp/constraints"
)

var (
	// Zero is constant value representing zero.
	Zero = newConst([]byte{0})
	// One is constant value representing one.
	One = newConst([]byte{1})
)

var _ Expr = Const{}

// Const is a constant expression of arbitrary width.
type Const struct {
	bs []byte
}

func newConst(bs []byte) Const {
	return Const{bs: bs}
}

// NewConst creates constant of width w out of bytes b. Bytes in b are encoded
// using little-endian encoding -> b[0] is the least significant byte of the
// value.
//
// Value of b is copied into an internal buffer. Consequently user is free to
// use b once call to this function is completed.
func NewConst(b []byte, w Width) Const {
	bCopy := make([]byte, w)
	copy(bCopy, b)
	return newConst(bCopy)
}

// NewConstUint converts any uint value into Const of width w.
//
// This method will panic in case val doesn't fit w bytes. It's allowed to
// convert val of type wider than w, but all bytes of val higher than w has to
// be zero bytes.
func NewConstUint[T constraints.Unsigned](val T, w Width) Const {
	valCopy := val

	bs := make([]byte, w)
	for i := range bs {
		bs[i] = byte(val)
		val >>= 8
	}

	if val > 0 {
		panic(fmt.Sprintf("value of type %T doesn't fit to value of width %d: %d",
			val, w, valCopy))
	}

	return newConst(bs)
}

// NewConstInt converts any int value into Const of width w.
//
// This method will panic in case val doesn't fit w bytes. It's allowed to
// convert val of type wider than w, but all bytes of val higher than w has to
// be sign extension of last bit of byte [w-1].
func NewConstInt[T constraints.Signed](val T, w Width) Const {
	valCopy := val

	bs := make([]byte, w)
	for i := range bs {
		bs[i] = byte(val)
		val >>= 8
	}

	if val != 0 && (val != -1 || bs[len(bs)-1] < 128) {
		panic(fmt.Sprintf("value of type %T doesn't fit to value of width %d: %d",
			val, w, valCopy))
	}

	return newConst(bs)
}

// ConstFromUint converts any unsigned integer into Const of the same width as T
// has.
func ConstFromUint[T constraints.Unsigned](val T) Const {
	return NewConstUint(val, Width(unsafe.Sizeof(val)))
}

// ConstFromInt converts any signed integer into Const of the same width as T
// has.
func ConstFromInt[T constraints.Signed](val T) Const {
	return NewConstInt(val, Width(unsafe.Sizeof(val)))
}

// Bytes returns array of bytes stored in c.
//
// The value returned must be treated as read-only. Any modification of returned
// value will modify c as well, which is prohibited given that expressions are
// immutable.
func (c Const) Bytes() []byte { return c.bs }

// Width returns width of c.
func (c Const) Width() Width { return Width(len(c.bs)) }

// Equal checks constant equality.
func (c1 Const) Equal(c2 Const) bool { return bytes.Equal(c1.bs, c2.bs) }
func (Const) internalExpr()          {}

// WithWidth returns new Const with same content as c but with width w.
func (c Const) WithWidth(w Width) Const {
	// Those 2 branches are performance optimization to avoid unnecessary
	// allocations.
	if c.Width() == w {
		return c
	} else if c.Width() > w {
		return newConst(c.bs[:w])
	}

	return NewConst(c.Bytes(), w)
}

func nonzeroUpperIdx(b []byte) int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0 {
			return i
		}
	}

	return 0
}

// ConstUint converts Const into an arbitrary uint type. The boolean return
// value indicates if Const value fits T. If Const value doesn't fit T, returned
// value contains lower sizeof(T) bytes.
//
// Stupid Go 1.18 doesn't allow type parameters for methods so this has to be a
// function...
func ConstUint[T constraints.Unsigned](c Const) (T, bool) {
	fits := true
	var val T

	idx := nonzeroUpperIdx(c.Bytes())
	if size := unsafe.Sizeof(val); uintptr(idx) >= size {
		idx, fits = int(size-1), false
	}

	for i := idx; i >= 0; i-- {
		val <<= 8
		val |= T(c.Bytes()[i])
	}

	return val, fits
}
