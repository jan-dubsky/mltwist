package expreval

import (
	"fmt"
	"math"
	"math/big"
	"mltwist/pkg/expr"
	"unsafe"
)

func init() {
	// For bit shifts, we assume that we are able to represent any possible
	// bit shift of register using uint64. To do so, possible width of
	// register shifted must be strictly less than 2^64 bits => 2^61 bytes.
	// As 2^64 is not representable as uint64 and we have more than enough
	// space, we have decided to use 2^63 bits as allowed maximum => 2^60
	// bytes.
	//
	// Given that 2^60 bytes is more than 1EB (exabytes), we can consider
	// this assumption to be reasonable. Current maximal widths of registers
	// are at most tenths of bytes (i.e. AVX512), so we have more than
	// enough spare capacity to grow to with this implementation.
	if unsafe.Sizeof(expr.Width(0))*8 > 60 {
		panic("precondition of this package is uint64 must be able to capture " +
			"any number of bits expr.Width can express")
	}
}

// Add calculates sum of val1 and val2 of width w.
func Add(val1 Value, val2 Value, w expr.Width) Value {
	sum := make([]byte, w)
	bytes1, bytes2 := val1.SetWidth(w).bytes(), val2.SetWidth(w).bytes()

	var carry bool
	for i := range sum {
		b1, b2 := bytes1[i], bytes2[i]
		res := b1 + b2

		var newCarry = res < b1 || (carry && res == math.MaxUint8)
		if carry {
			res += 1
		}

		sum[i] = res
		carry = newCarry
	}

	return newValue(sum)
}

// Sub calculates difference of val1 and val2 of width w.
func Sub(val1 Value, val2 Value, w expr.Width) Value {
	diff := make([]byte, w)
	bytes1, bytes2 := val1.SetWidth(w).bytes(), val2.SetWidth(w).bytes()

	var carry bool
	for i := range diff {
		b1, b2 := bytes1[i], bytes2[i]
		res := b1 - b2

		var newCarry = res > b1 || (carry && res == 0)
		if carry {
			res -= 1
		}

		diff[i] = res
		carry = newCarry
	}

	return newValue(diff)
}

func shiftUint64(v Value, w expr.Width) (uint64, uint8, bool) {
	vInt := v.bigInt(w)

	// Is greater than any allowed width of value -> The whole shifted value
	// will be shifted away -> result is always zero independently on the
	// value.
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

func bitLsh(bs []byte, shift uint8) {
	if shift >= 8 || shift == 0 {
		panic(fmt.Sprintf("invalid bit shift: %d", shift))
	}

	antiShift := 8 - shift
	bs[len(bs)-1] = bs[len(bs)-1] << shift
	for i := len(bs) - 2; i >= 0; i-- {
		bs[i+1] |= bs[i] >> antiShift
		bs[i] <<= shift
	}
}

// Lsh left shifts value val1 of val2 bits. Bottom-most val2 bits are set to
// zeros.
func Lsh(val1 Value, val2 Value, w expr.Width) Value {
	byteShift, bitShift, ok := shiftUint64(val2, w)
	if !ok {
		return Value{}.SetWidth(w)
	}

	val1Ext := val1.SetWidth(w).bytes()
	shifted := make([]byte, w)
	for i := 0; i < int(w-expr.Width(byteShift)); i++ {
		shifted[i+int(byteShift)] = val1Ext[i]
	}

	if bitShift != 0 {
		bitLsh(shifted, bitShift)
	}

	return newValue(shifted)
}

func bitRsh(bs []byte, shift uint8) {
	if shift >= 8 || shift == 0 {
		panic(fmt.Sprintf("invalid bit shift: %d", shift))
	}

	antiShift := 8 - shift
	bs[0] = bs[0] >> shift
	for i := 1; i < len(bs); i++ {
		bs[i-1] |= bs[i] << antiShift
		bs[i] >>= shift
	}
}

// Rsh right shifts value val1 of val2 bits. Top-most val2 bits are filled to
// zeros.
func Rsh(val1 Value, val2 Value, w expr.Width) Value {
	byteShift, bitShift, ok := shiftUint64(val2, w)
	if !ok {
		return Value{}.SetWidth(w)
	}

	bytes1 := val1.SetWidth(w).bytes()
	shiftedLen := int(w - expr.Width(byteShift))
	shifted := make([]byte, w)
	for i := 0; i < shiftedLen; i++ {
		shifted[i] = bytes1[i+int(byteShift)]
	}

	if bitShift != 0 {
		bitRsh(shifted[:shiftedLen], bitShift)
	}
	return newValue(shifted)
}

// Mul unsigned multiplies val1 and val2. The result has width w.
func Mul(val1 Value, val2 Value, w expr.Width) Value {
	val1Int, val2Int := val1.bigInt(w), val2.bigInt(w)
	product := (&big.Int{}).Mul(val1Int, val2Int)
	return parseBigInt(product).SetWidth(w)
}

// Div calculates val1 (unsigned) divided by div2. The produced value is of
// width w. If val2 is zero, this method returns all ones of width w.
func Div(val1 Value, val2 Value, w expr.Width) Value {
	val2Int := val2.bigInt(w)

	// Special-case division by zero.
	if val2Int.Cmp(&big.Int{}) == 0 {
		bytes := make([]byte, w)
		for i := range bytes {
			bytes[i] = math.MaxUint8
		}
		return newValue(bytes)
	}

	div := (&big.Int{}).Div(val1.bigInt(w), val2Int)
	return parseBigInt(div).SetWidth(w)
}

func Nand(val1 Value, val2 Value, w expr.Width) Value {
	result := make([]byte, w)
	bytes1, bytes2 := val1.SetWidth(w).bytes(), val2.SetWidth(w).bytes()

	for i := range result {
		result[i] = ^(bytes1[i] & bytes2[i])
	}

	return newValue(result)
}
