package expr

import (
	"decomp/pkg/expr/internal/value"
	"math"
)

type AddI struct {
	width Width
	ex1   Expr
	ex2   Expr
}

func NewAddI(width Width, ex1 Expr, ex2 Expr) AddI {
	return AddI{
		width: width,
		ex1:   ex1,
		ex2:   ex2,
	}
}

func (a AddI) value() (value.Value, bool) {
	v1, ok := a.ex1.value()
	if !ok {
		return value.Value{}, false
	}

	v2, ok := a.ex2.value()
	if !ok {
		return value.Value{}, false
	}

	v, _ := addI(a.width, v1, v2)
	return v, true
}

func maxByte(b1 byte, b2 byte) byte {
	if b1 > b2 {
		return b1
	} else {
		return b2
	}
}

func addI(width Width, v1 value.Value, v2 value.Value) (value.Value, bool) {
	result := make([]byte, width)
	var overflow byte

	for i := 0; i < len(result); i++ {
		sum := v1.Index(i) + v2.Index(i)
		result[i] = sum + overflow

		overflow = 0
		if (sum < maxByte(v1.Index(i), v2.Index(i))) ||
			(sum == math.MaxInt8 && overflow == 1) {
			overflow = 1
		}
	}

	return value.New(result), overflow == 1
}
