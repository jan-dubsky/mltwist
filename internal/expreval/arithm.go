package expreval

import (
	"decomp/pkg/expr"
	"fmt"
	"math"
	"math/big"
	"unsafe"
)

func init() {
	// For bit shifts, we assume that we are able to represent any possible
	// bit shift of register using uint64. To do so, possible width of
	// register shifted must be less than 2^64 bits => 2^61 bytes.
	//
	// Given that 2^61 bytes is more than 2EB (exabytes), we can consider
	// this assumption to be reasonable. Current maximal widths of registers
	// are at most tenths of bytes (i.e. AVX512), so we have more than
	// enough spare capacity to grow to with this implementation.
	if (unsafe.Sizeof(expr.Width(0)) * 8) > 60 {
		panic("precondition of this package is uint64 must be able to capture " +
			"any number of bits expr.Width can express")
	}
}

func add(val1 value, val2 value, w expr.Width) value {
	sum := make(value, w)
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)

	var carry bool
	for i := range sum {
		v1, v2 := val1Ext[i], val2Ext[i]
		res := v1 + v2

		var newCarry = res < v1 || (carry && res == math.MaxUint8)
		if carry {
			res += 1
		}

		sum[i] = res
		carry = newCarry
	}

	return sum
}

func sub(val1 value, val2 value, w expr.Width) value {
	diff := make(value, w)
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)

	var carry bool
	for i := range diff {
		v1, v2 := val1Ext[i], val2Ext[i]
		res := v1 - v2

		var newCarry = res > v1 || (carry && res == 0)
		if carry {
			res -= 1
		}

		diff[i] = res
		carry = newCarry
	}

	return diff
}

func shiftUint64(v value, w expr.Width) (uint64, uint8, bool) {
	vInt := v.bigInt(w)
	if !vInt.IsUint64() {
		return 0, 0, false
	}

	rawShift := vInt.Uint64()

	byteShift := rawShift / 8
	if byteShift >= uint64(w) {
		return 0, 0, false
	}

	return byteShift, uint8(rawShift % 8), true
}

func bitLsh(val value, shift uint8) {
	if shift >= 8 || shift == 0 {
		panic(fmt.Sprintf("invalid bit shift: %d", shift))
	}

	antiShift := 8 - shift
	val[len(val)-1] = val[len(val)-1] << shift
	for i := len(val) - 2; i >= 0; i-- {
		val[i+1] |= val[i] >> antiShift
		val[i] <<= shift
	}
}

func lsh(val1 value, val2 value, w expr.Width) value {
	byteShift, bitShift, ok := shiftUint64(val2, w)
	if !ok {
		return value{}.setWidth(w)
	}

	val1Ext := val1.setWidth(w)
	shifted := make(value, w)
	for i := 0; i < int(w-expr.Width(byteShift)); i++ {
		shifted[i+int(byteShift)] = val1Ext[i]
	}

	if bitShift != 0 {
		bitLsh(shifted, bitShift)
	}

	return shifted
}

func bitRsh(val value, shift uint8) {
	if shift >= 8 || shift == 0 {
		panic(fmt.Sprintf("invalid bit shift: %d", shift))
	}

	antiShift := 8 - shift
	val[0] = val[0] >> shift
	for i := 1; i < len(val); i++ {
		val[i-1] |= val[i] << antiShift
		val[i] >>= shift
	}
}

func rsh(val1 value, val2 value, w expr.Width) value {
	byteShift, bitShift, ok := shiftUint64(val2, w)
	if !ok {
		return value{}.setWidth(w)
	}

	val1Ext := val1.setWidth(w)
	shifted := make(value, w-expr.Width(byteShift))
	for i := range shifted {
		shifted[i] = val1Ext[i+int(byteShift)]
	}

	if bitShift != 0 {
		bitRsh(shifted, bitShift)
	}
	return shifted.setWidth(w)
}

func mul(val1 value, val2 value, w expr.Width) value {
	val1Int, val2Int := val1.bigInt(w), val2.bigInt(w)
	product := (&big.Int{}).Mul(val1Int, val2Int)
	return valFromBigInt(product).setWidth(w)
}

func div(val1 value, val2 value, w expr.Width) value {
	val2Int := val2.bigInt(w)

	// Special-case division by zero.
	if val2Int.Cmp(&big.Int{}) == 0 {
		div := make(value, w)
		for i := range div {
			div[i] = math.MaxUint8
		}
		return div
	}

	div := (&big.Int{}).Div(val1.bigInt(w), val2Int)
	return valFromBigInt(div).setWidth(w)
}

func mod(val1 value, val2 value, w expr.Width) value {
	val2Int := val2.bigInt(w)

	// Special-case division by zero.
	if val2Int.Cmp(&big.Int{}) == 0 {
		return val1.setWidth(w)
	}

	mod := (&big.Int{}).Mod(val1.bigInt(w), val2Int)
	return valFromBigInt(mod).setWidth(w)
}

func bitOp(
	val1 value,
	val2 value,
	w expr.Width,
	byteFunc func(v1 byte, v2 byte) byte,
) value {
	result := make(value, w)
	val1Ext, val2Ext := val1.setWidth(w), val2.setWidth(w)

	for i := range result {
		result[i] = byteFunc(val1Ext[i], val2Ext[i])
	}

	return result
}

func and(val1 value, val2 value, w expr.Width) value {
	return bitOp(val1, val2, w, func(v1, v2 byte) byte { return v1 & v2 })
}

func or(val1 value, val2 value, w expr.Width) value {
	return bitOp(val1, val2, w, func(v1, v2 byte) byte { return v1 | v2 })
}

func xor(val1 value, val2 value, w expr.Width) value {
	return bitOp(val1, val2, w, func(v1, v2 byte) byte { return v1 ^ v2 })
}
