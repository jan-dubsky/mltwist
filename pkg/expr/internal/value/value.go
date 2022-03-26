package value

import "math"

// Value represents an arbitrary integer or floating point value.
//
// Integer values are encoded as little endian value using one's complement.
//
// Floating point are represented using IEEE 754 standard.
type Value struct {
	b []byte
}

func New(b []byte) Value {
	v := new(b)
	v.Compact()
	return v
}

func new(b []byte) Value {
	return Value{b: b}
}

func (v Value) Len() int { return len(v.b) }

func (v Value) Index(i int) byte {
	if i >= len(v.b) {
		return 0
	}

	return v.b[i]
}

func (v Value) Bytes() []byte { return v.b }

func (v *Value) Compact() {
	for i := len(v.b) - 1; i > 0; i-- {
		if v.b[i] != 0 {
			break
		}

		v.b = v.b[:len(v.b)-1]
	}
}

func (v Value) Clone() Value {
	b := make([]byte, len(v.b))
	copy(b, v.b)
	return new(b)
}

func (v Value) Zero() bool {
	for _, b := range v.b {
		if b != 0 {
			return false
		}
	}
	return true
}

func (v Value) Negative() bool { return v.b[len(v.b)-1]&128 != 0 }

func invertBits(bs []byte) {
	for i, b := range bs {
		bs[i] = ^b
	}
}

func (v Value) Negate() {
	if v.Negative() {
		invertBits(v.b)

		for i, b := range v.b {
			if b != math.MaxUint8 {
				v.b[i] = b + 1
				break
			}
		}
	} else if !v.Zero() {
		for i, b := range v.b {
			if b != 0 {
				v.b[i] = b - 1
				break
			}
		}

		invertBits(v.b)
	}
}

func Equal(v1 Value, v2 Value) bool {
	if v1.Len() != v2.Len() {
		return false
	}

	for i := range v1.b {
		if v1.b[i] != v2.b[i] {
			return false
		}
	}

	return true
}
