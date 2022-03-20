package riscv

import "fmt"

// immType represents one of 6 potential RISC V encodings of an immediate value
// in an instruction opcode.
//
// Allowed values of immType are then 0-5 inclusive, where zero value represents
// an opcode layout without an immediate value. The fact that zero value is used
// to express such layout is convenient as it allows us not to specify immediate
// layour for instruction opcodes without immediate value and laverage compiler
// zero value initialization instead.
type immType uint8

const (
	// immTypeR is an instruction opcode format which has no immediate value
	// encoded in it.
	//
	// As immTypeR contains no immediate value at all, it's convenient to
	// express it using zero value. This allows not to specify immediate at
	// all for opcodes not encoding an immediate value which increases
	// readability of instruction opcode readability.
	immTypeR immType = iota

	// immTypeI is an instruction opcode format which encodes lower 12 bits
	// of an immediate value in top 12 bits if the opcode.
	//
	// Immediate bits parsed from an opcode are sign extended to fill
	// register width.
	immTypeI

	// immTypeS is an instruction opcode format which encodes lowest 12 bits
	// of an immediate value in bits [7:11] and [25:31] of an instruction
	// opcode.
	//
	// Immediate bits parsed from an opcode are sign extended to fill
	// register width.
	immTypeS

	// immTypeB is an instruction opcode format which encodes 12 bits of an
	// immediate value in instruction opcode bits: [8:11], [25:30], [7], and
	// [31] respectively.
	//
	// Immediate bits parsed from opcode bits are then understood as if
	// those were bottom 12 bits of a value which would be then logically
	// shifted to left by one and sign extended to fill register width.
	//
	// The bit order used to represent an immediate value feels to be weird
	// in the first, second and even third sight. But authors of the
	// specification claim that this is simpler to implement in hardware. Or
	// more specifically that this format will cause slight signal delay
	// compared to using bits in order, but that compared to that solution
	// this design requires less gates to be used and decreases the system
	// complexity. So why not to trust them.
	//
	// Bit 31 of the opcode still represents sign bit. According to authors
	// of the specification this allows to decrease hardware complexity as
	// in all opcodes the 31 bit is the sign bit.
	immTypeB

	// immTypeU is an instruction opcode format which encodes upper 20 bits
	// on an immediate value (bits [12:31]) in bits [12:31] of an
	// instruction opcode.
	//
	// Bottom 12 bits of an immediate value are unset - filled with zeros.
	immTypeU

	// immTypeJ is an instruction opcode format which encodes 12 bits of an
	// immediate value in instruction opcode bits: [21:30], [20], [12:19],
	// and [31] respectively.
	//
	// Immediate bits parsed from opcode bits are then understood as if
	// those were bottom 20 bits of a value which would be then logically
	// shifted by 1 and sign extended to fill the register width.
	//
	// Even though this format seems to be super crazy (even more than
	// immTypeB format), the authors state that it has been designed in this
	// way to be the most similar to other instruction formats as possible.
	// Authors state that this similarity then allows to share gates with
	// other instruction types and to simplify the hardware this way.
	//
	// Bit 31 of the opcode still represents sign bit. According to authors
	// of the specification this allows to decrease hardware complexity as
	// in all opcodes the 31 bit is the sign bit.
	immTypeJ
)

// parseBitRange parses bits in range [begin,end) in value into lowest bytes if
// a return value.
func parseBitRange(value uint32, begin uint8, end uint8) uint32 {
	endMask := (uint32(1) << end) - 1
	return (value & endMask) >> begin
}

// signExtend perform integer sign extentions on unsigned where value of bit
// with index signBitIdx is copied to all higher bits.
//
// As sign extension of 32bit value doesn't make sense for values which are
// already 32bit width, this method will panic if either unsigned has bit 31
// set, or if signBitIdx >= 31. This method will also panic of any bit of
// unsigned with higher index than signBitIdx is set, as in such a case
// sign-extension makes no sense.
func signExtend(unsigned uint32, signBitIdx uint8) int32 {
	if unsigned >= uint32(1)<<31 {
		panic(fmt.Sprintf(
			"invalid unsigned value for sign extension in bit %d: 0x%x",
			signBitIdx, unsigned))
	}
	if signBitIdx >= 31 {
		panic(fmt.Sprintf("invalid sign bit index: %d", signBitIdx))
	}

	signBitMask := uint32(1) << signBitIdx
	if max := signBitMask << 1; unsigned >= max {
		panic(fmt.Sprintf("bit above index %d is set: 0x%x", signBitIdx, unsigned))
	}

	if unsigned&signBitMask == 0 {
		return int32(unsigned)
	}

	mask := -(int32(1) << (signBitIdx + 1))
	return int32(unsigned) | mask
}

// parseValue parses an immediate value out of an instruction. Based on an
// immediate type, this method will return either (0, false) for immediate type
// R and an immediate value and true for any other type. This method will panic
// for unknown immediate type.
func (t immType) parseValue(value uint32) (int32, bool) {
	switch t {
	case immTypeR:
		return 0, false
	case immTypeI:
		val := signExtend(parseBitRange(value, 20, 32), 11)
		return val, true
	case immTypeS:
		low := parseBitRange(value, 7, 12)
		high := parseBitRange(value, 25, 32)
		val := signExtend((high)<<5|low, 11)
		return val, true
	case immTypeB:
		first := parseBitRange(value, 8, 12)
		second := parseBitRange(value, 25, 31)
		third := parseBitRange(value, 7, 8)
		sign := parseBitRange(value, 31, 32)

		// Bit [0] is set to 0 be hardware.
		unsigned := (first << 1) | (second << 5) | (third << 11) | (sign << 12)
		val := signExtend(unsigned, 12)
		return val, true
	case immTypeU:
		unsigned := parseBitRange(value, 12, 32)
		return int32(unsigned << 12), true
	case immTypeJ:
		first := parseBitRange(value, 21, 31)
		second := parseBitRange(value, 20, 21)
		third := parseBitRange(value, 12, 20)
		sign := parseBitRange(value, 31, 32)

		// Bit [0] is set to 0 be hardware.
		unsigned := (first << 1) | (second << 11) | (third << 12) | (sign << 20)
		val := signExtend(unsigned, 20)
		return val, true
	default:
		panic(fmt.Sprintf("unknown immediate type: %v", t))
	}
}
