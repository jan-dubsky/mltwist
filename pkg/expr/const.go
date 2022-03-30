package expr

import (
	"math/big"
)

var (
	Zero = NewConst(big.NewInt(0))
	One  = NewConst(big.NewInt(1))
)

type Const struct {
	v *big.Int
}

func NewConst(v *big.Int) Const {
	return Const{v: v}
}

func NewAddrConst(a uint64) Const {
	var v big.Int
	return NewConst(v.SetUint64(a))
}
