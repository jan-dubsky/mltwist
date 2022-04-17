package expr

import (
	"fmt"
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

// NewConstUint converts any uint value into Const of width w.
//
// This method will panic in case val doesn't fit w bytes. It's allowed to
// convert val of type wider than w, but all bytes of val higher than w has to
// be zero bytes.
func NewConstUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint](val T, w Width) Const {
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
func NewConstInt[T ~int8 | ~int16 | ~int32 | ~int64](val T, w Width) Const {
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
func (Const) internalExpr()   {}
