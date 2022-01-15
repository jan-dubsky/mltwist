package opcode

func applyMask(bytes []byte, mask []byte) []byte {
	masked := make([]byte, len(mask))
	for i := range mask {
		masked[i] = bytes[i] & mask[i]
	}
	return masked
}

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

func byteLE(first []byte, second []byte) bool {
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

	return true
}
