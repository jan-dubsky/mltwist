package expreval

import "mltwist/pkg/expr"

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

// Ltu compares w bytes of val1 and val2 and checks if val1 is less then val2
// using unsigned integer comparison.
func Ltu(val1 Value, val2 Value, w expr.Width) bool {
	val1Ext, val2Ext := val1.SetWidth(w), val2.SetWidth(w)
	return ltu(val1Ext, val2Ext)
}
