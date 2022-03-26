package expr

import "decomp/pkg/expr/internal/value"

var _ Expr = Shift{}

type Shift struct {
	width Width
	val   Expr
	shift Expr
}

// TODO: Doc-comment (negative shift == >>, positive shift == <<).
func NewShift(width Width, val Expr, shift Expr) Shift {
	return Shift{
		width: width,
		val:   val,
		shift: shift,
	}
}

func (s Shift) value() (value.Value, bool) {
	v, ok := s.val.value()
	if !ok {
		return value.Value{}, false
	}

	shift, ok := s.shift.value()
	if !ok {
		return value.Value{}, false
	}

	if shift.Zero() {
		return v, true
	}

	panic("not yet implemented")

	//if shift.Negative() {
	//
	//} else {
	//
	//}
}

//func rshift(width uint8, v Expr, shift uint64) []byte {
//	result := make([]byte, width)
//
//}
