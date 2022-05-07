package expreval

import "mltwist/pkg/expr"

func eq(val1 Value, val2 Value) bool {
	for i := range val1.bytes() {
		if val1.bytes()[i] != val2.bytes()[i] {
			return false
		}
	}

	return true
}

func Eq(val1 Value, val2 Value, w expr.Width) bool {
	val1Ext, val2Ext := val1.SetWidth(w), val2.SetWidth(w)
	return eq(val1Ext, val2Ext)
}

func ltu(val1 Value, val2 Value) bool {
	bytes1, bytes2 := val1.bytes(), val2.bytes()
	for i := len(bytes1) - 1; i >= 0; i-- {
		if bytes1[i] < bytes2[i] {
			return true
		} else if bytes1[i] > bytes2[i] {
			return false
		}
	}

	return false
}

func Ltu(val1 Value, val2 Value, w expr.Width) bool {
	val1Ext, val2Ext := val1.SetWidth(w), val2.SetWidth(w)
	return ltu(val1Ext, val2Ext)
}

func Leu(val1 Value, val2 Value, w expr.Width) bool {
	val1Ext, val2Ext := val1.SetWidth(w), val2.SetWidth(w)
	return ltu(val1Ext, val2Ext) || eq(val1Ext, val2Ext)
}

func negative(v Value) bool {
	bs := v.bytes()
	return bs[len(bs)-1]&0x80 != 0
}

func lts(val1 Value, val2 Value) bool {
	val1Neg, val2Neg := negative(val1), negative(val2)
	if val1Neg != val2Neg {
		return val1Neg
	}

	if val1Neg {
		return ltu(val2, val1)
	} else {
		return ltu(val1, val2)
	}
}

func Lts(val1 Value, val2 Value, w expr.Width) bool {
	val1Ext, val2Ext := val1.SetWidth(w), val2.SetWidth(w)
	return lts(val1Ext, val2Ext)
}

func Les(val1 Value, val2 Value, w expr.Width) bool {
	val1Ext, val2Ext := val1.SetWidth(w), val2.SetWidth(w)
	return lts(val1Ext, val2Ext) || eq(val1Ext, val2Ext)
}
