package expreval

import (
	"decomp/pkg/expr"
	"math/big"
)

type Value []byte

func ParseConst(e expr.Const) Value { return Value(e.Bytes()) }

func parseBigInt(i *big.Int) Value {
	v := Value(i.Bytes())
	revertBytes(v)
	return v
}

func (v Value) SetWidth(w expr.Width) Value {
	if int(w) <= len(v) {
		return v[:w]
	}

	extended := make([]byte, w)
	copy(extended, v)
	return extended
}

func revertBytes(v Value) {
	for i := 0; i < len(v)/2; i++ {
		v[i], v[len(v)-1-i] = v[len(v)-1-i], v[i]
	}
}

func (v Value) clone() Value {
	val := make(Value, len(v))
	for i, b := range v {
		val[i] = b
	}
	return val
}

func (v Value) bigInt(w expr.Width) *big.Int {
	vCut := v
	if int(w) < len(v) {
		vCut = v.SetWidth(w)
	}

	vBig := vCut.clone()
	revertBytes(vBig)

	return (&big.Int{}).SetBytes(vBig)
}

func (v Value) Const(w expr.Width) expr.Const {
	return expr.NewConst(v, expr.Width(len(v)))
}
