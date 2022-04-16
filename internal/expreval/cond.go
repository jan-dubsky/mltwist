package expreval

import "decomp/pkg/expr"

func eqInternal(val1 value, val2 value) bool {
	for i := range val1 {
		if val1[i] != val2[i] {
			return false
		}
	}

	return true
}

func eq(val1 value, val2 value, w expr.Width) bool {
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)
	return eqInternal(val1Ext, val2Ext)
}

func ltuInternal(val1 value, val2 value) bool {
	for i := len(val1) - 1; i >= 0; i-- {
		if val1[i] < val2[i] {
			return true
		} else if val1[i] > val2[i] {
			return false
		}
	}

	return false
}

func ltu(val1 value, val2 value, w expr.Width) bool {
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)
	return ltuInternal(val1Ext, val2Ext)
}

func leu(val1 value, val2 value, w expr.Width) bool {
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)
	return ltuInternal(val1Ext, val2Ext) || eqInternal(val1Ext, val2Ext)
}

func negative(v value) bool { return v[len(v)-1]&0x80 != 0 }

func ltsInternal(val1 value, val2 value) bool {
	val1Neg, val2Neg := negative(val1), negative(val2)
	if val1Neg != val2Neg {
		return val1Neg
	}

	if val1Neg {
		return ltuInternal(val2, val1)
	} else {
		return ltuInternal(val1, val2)
	}
}

func lts(val1 value, val2 value, w expr.Width) bool {
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)
	return ltsInternal(val1Ext, val2Ext)
}

func les(val1 value, val2 value, w expr.Width) bool {
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)
	return ltsInternal(val1Ext, val2Ext) || eqInternal(val1Ext, val2Ext)
}
