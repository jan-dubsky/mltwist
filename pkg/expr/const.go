package expr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

var (
	// Zero is constant value representing zero.
	Zero = newConst([]byte{0})
	// One is constant value representing one.
	One = newConst([]byte{1})
)

var _ Expr = Const{}

type Const struct {
	b []byte
}

func newConst(b []byte) Const {
	return Const{
		b: b,
	}
}

// NewConst creates constant of width w out of bytes b.
//
// Value of b is copied into an internal buffer. Consequently user is free to
// use b once call to this function is completed.
func NewConst(b []byte, w Width) Const {
	bCopy := make([]byte, w)
	copy(bCopy, b)
	return newConst(bCopy)
}

type uints interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

type ints interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// NewConstUint converts any uint value into Const of width w.
//
// This method will panic in case val doesn't fit w bytes. It's allowed to
// convert val of type wider than w, but all bytes of val higher than w has to
// be zero bytes.
func NewConstUint[T uints](val T, w Width) Const {
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
func NewConstInt[T ints](val T, w Width) Const {
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

func (c Const) Bytes() []byte { return c.b }
func (c Const) Width() Width  { return Width(len(c.b)) }

// Equal checks constant equality.
func (c1 Const) Equal(c2 Const) bool { return bytes.Equal(c1.b, c2.b) }
func (Const) internalExpr()          {}

func nonzeroUpperIdx(b []byte) int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0 {
			return i
		}
	}

	return 0
}

// ConstUint converts Const into an arbitrary uint type. The boolean return
// value indicates if conversion was successful or Const value doesn't fit T. In
// the latter case, the value of returned uint is undefined.
func ConstUint[T uints](c Const) (T, bool) {
	var dummy T
	TSize := unsafe.Sizeof(dummy)

	if uintptr(c.Width()) > TSize {
		return 0, false
	}

	idx := int(c.Width())
	if uintptr(idx) > TSize {
		idx = nonzeroUpperIdx(c.Bytes())
		if uintptr(idx) > TSize {
			return 0, false
		}
	}

	val := binary.LittleEndian.Uint64(c.Bytes())
	return T(val), true
}
