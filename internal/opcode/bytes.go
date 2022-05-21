package opcode

// applyMask sets to zero bits which are not included in the mask. Both bytes
// and mask must have the same length otherwise this function might panic.
func applyMask(bytes []byte, mask []byte) []byte {
	masked := make([]byte, len(mask))
	for i := range mask {
		masked[i] = bytes[i] & mask[i]
	}
	return masked
}

// byteEQ compares if first and second are equal.
func byteEQ(first []byte, second []byte) bool {
	if len(first) != len(second) {
		return false
	}

	for i := range first {
		if first[i] != second[i] {
			return false
		}
	}

	return true
}

// byteLT compares first byte array to the second one and returns true of either
// first is shorter and second or both first and second have equal length, but
// the big endian of first value is less then big endian value of second. This
// method defines ordering for byte arrays.
func byteLT(first []byte, second []byte) bool {
	if len(first) < len(second) {
		return true
	} else if len(first) > len(second) {
		return false
	}

	for i := range first {
		if first[i] < second[i] {
			return true
		} else if first[i] > second[i] {
			return false
		}
	}

	return false
}
