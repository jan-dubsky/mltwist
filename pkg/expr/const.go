package expr

import (
	"math/big"
)

var (
	Zero = NewConst(big.NewInt(0), Width8)
	One  = NewConst(big.NewInt(1), Width8)
)

var _ Expr = Const{}

type Const struct {
	v *big.Int
	w Width
}

func NewConst(v *big.Int, w Width) Const { return Const{v: v} }

func NewConstUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64](val T, w Width) Const {
	var v big.Int
	return NewConst(v.SetUint64(uint64(val)), w)
}

func NewConstInt[T ~int8 | ~int16 | ~int32 | ~int64](val T, w Width) Const {
	return NewConst(big.NewInt(int64(val)), w)
}

func (c Const) Width() Width { return c.w }
func (Const) internal()      {}
