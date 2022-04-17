package expreval

import "decomp/pkg/expr"

func eq(val1 Value, val2 Value) bool {
	for i := range val1 {
		if val1[i] != val2[i] {
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
	for i := len(val1) - 1; i >= 0; i-- {
		if val1[i] < val2[i] {
			return true
		} else if val1[i] > val2[i] {
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

func negative(v Value) bool { return v[len(v)-1]&0x80 != 0 }

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
