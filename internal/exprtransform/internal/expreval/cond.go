package expreval

import "mltwist/pkg/expr"

// Ltu compares w bytes of val1 and val2 and checks if val1 is less then val2
// using unsigned integer comparison.
func Ltu(val1 Value, val2 Value, w expr.Width) bool {
	bytes1, bytes2 := val1.setWidth(w).bytes(), val2.setWidth(w).bytes()
	for i := len(bytes1) - 1; i >= 0; i-- {
		if bytes1[i] < bytes2[i] {
			return true
		} else if bytes1[i] > bytes2[i] {
			return false
		}
	}

	// Both are equal.
	return false
}
