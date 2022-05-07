package expreval

import (
	"fmt"
	"math/big"
	"mltwist/pkg/expr"
)

// Value represents a number composet of arbitrary number of bytes.
type Value struct {
	bs []byte
}

// NewValue creates a new value comprising of bytes in b. Value overtakes
// ownership of b. Consequently the caller is not allowed to modify b after this
// call.
func NewValue(b []byte) Value {
	if ln := len(b); int(expr.Width(ln)) != ln {
		panic(fmt.Sprintf("byte array is too long: %d", ln))
	}

	return newValue(b)
}

func newValue(b []byte) Value { return Value{bs: b} }

// ParseConst converts contant expression into a value.
func ParseConst(e expr.Const) Value {
	// This implies that any modification of Value modified e as well. Given
	// that we never modify Value in this package this is safe performance
	// optimization.
	return newValue(e.Bytes())
}

// parseBigInt converts bit.Int into a value.
func parseBigInt(i *big.Int) Value {
	v := NewValue(i.Bytes())
	v.revertBytes()
	return v
}

// bytes return byte array stored in v.
func (v Value) bytes() []byte { return v.bs }

// width returns number of bytes in v as expr.Width.
func (v Value) width() expr.Width { return expr.Width(len(v.bytes())) }

// SetWidth creates a new value with width w.
func (v Value) SetWidth(w expr.Width) Value {
	if w <= v.width() {
		return newValue(v.bytes()[:w])
	}

	extended := make([]byte, w)
	copy(extended, v.bytes())
	return newValue(extended)
}

// revertBytes converts bytes from big to little endian order and vice versa.
//
// WARNING: This method modified v.
func (v *Value) revertBytes() {
	ln := len(v.bs)
	for i := 0; i < ln/2; i++ {
		v.bs[i], v.bs[ln-1-i] = v.bs[ln-1-i], v.bs[i]
	}
}

// clone creates a new copy of value.
func (v Value) clone() Value {
	b := make([]byte, v.width())
	copy(b, v.bytes())
	return NewValue(b)
}

func (v Value) bigInt(w expr.Width) *big.Int {
	// If w >= len(v), setting v into bigInt and setting v.SetWidth() is
	// equivalent. Not calling SetWidth is juts a performance optimization
	// not to copy the byte array multiple times.
	vCut := v
	if w < v.width() {
		vCut = v.SetWidth(w)
	}

	vBig := vCut.clone()
	vBig.revertBytes()

	return (&big.Int{}).SetBytes(vBig.bytes())
}

// Const returns new expr.Const containing value of v. The returned Const has
// width based on width of w.
func (v Value) Const(w expr.Width) expr.Const {
	return expr.NewConst(v.bytes(), v.width())
}
