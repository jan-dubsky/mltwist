package expreval

import (
	"math/big"
	"mltwist/pkg/expr"
)

// Value represents a number composet of arbitrary number of bytes.
type Value struct {
	bs []byte
}

func newValue(b []byte) Value { return Value{bs: b} }

// ParseConst converts constant expression into a value.
func ParseConst(e expr.Const) Value {
	// This implies that any modification of Value modified e as well. Given
	// that we never modify Value in this package this is safe performance
	// optimization.
	return newValue(e.Bytes())
}

// parseBigInt converts bit.Int into a value.
//
// WANIRNG: The value of i is made invalid by this method and it's not allowed
// to access i after this method returns.
func parseBigInt(i *big.Int, w expr.Width) Value {
	var bs []byte
	if ln := len(i.Bytes()); ln >= int(w) {
		// The byte order is big endian, so least significant bytes are
		// at the end of the array.
		bs = i.Bytes()[ln-int(w) : ln]
	} else {
		bs = make([]byte, w)
		i.FillBytes(bs)
	}

	revertBytes(bs)
	return newValue(bs)
}

// bytes return byte array stored in v.
func (v Value) bytes() []byte { return v.bs }

// width returns number of bytes in v as expr.Width.
func (v Value) width() expr.Width { return expr.Width(len(v.bytes())) }

// setWidth creates a new value with width w.
func (v Value) setWidth(w expr.Width) Value {
	if w <= v.width() {
		return newValue(v.bytes()[:w])
	}

	extended := make([]byte, w)
	copy(extended, v.bytes())
	return newValue(extended)
}

// revertBytes converts bytes from big to little endian order and vice versa.
//
// WARNING: This method modified bs.
func revertBytes(bs []byte) {
	ln := len(bs)
	for i := 0; i < ln/2; i++ {
		bs[i], bs[ln-1-i] = bs[ln-1-i], bs[i]
	}
}

// clone creates a new copy of value.
func (v Value) clone() Value {
	b := make([]byte, v.width())
	copy(b, v.bytes())
	return newValue(b)
}

func (v Value) bigInt(w expr.Width) *big.Int {
	// If w >= len(v), setting v into bigInt and setting v.SetWidth() is
	// equivalent. Not calling SetWidth is just a performance optimization
	// not to copy the byte array multiple times.
	vCut := v
	if w < v.width() {
		vCut = v.setWidth(w)
	}

	vBig := vCut.clone()
	revertBytes(vBig.bs)

	return (&big.Int{}).SetBytes(vBig.bytes())
}

// Const returns new expr.Const containing value of v. The returned Const has
// width based on width of w.
func (v Value) Const(w expr.Width) expr.Const {
	return expr.NewConst(v.bytes(), v.width())
}
