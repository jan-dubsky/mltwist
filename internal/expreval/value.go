package expreval

import (
	"decomp/pkg/expr"
	"math/big"
)

type value []byte

func valFromBigInt(i *big.Int) value {
	v := value(i.Bytes())
	revertBytes(v)
	return v
}

func (v value) setWidth(w expr.Width) value {
	if int(w) <= len(v) {
		return v[:w]
	}

	extended := make([]byte, w)
	for i, b := range v {
		extended[i] = b
	}

	return extended
}

func revertBytes(v value) {
	for i := 0; i < len(v)/2; i++ {
		v[i], v[len(v)-1-i] = v[len(v)-1-i], v[i]
	}
}

func (v value) clone() value {
	val := make(value, len(v))
	for i, b := range v {
		val[i] = b
	}
	return val
}

func (v value) bigInt(w expr.Width) *big.Int {
	vCut := v
	if int(w) < len(v) {
		vCut = v.setWidth(w)
	}

	vBig := vCut.clone()
	revertBytes(vBig)

	return (&big.Int{}).SetBytes(vBig)
}
