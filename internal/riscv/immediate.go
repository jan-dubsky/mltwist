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

func parseBitRange(b InstrBytes, begin uint8, end uint8) uint32 {
	value := b.uint32()
	endMask := (uint32(1) << end) - 1
	return (value & endMask) >> begin
}

func signExtend(unsigned uint32, signBit uint8) int32 {
	if unsigned >= uint32(1)<<31 {
		panic(fmt.Sprintf(
			"invalid unsigned value for sign extension in bit %d: 0x%x",
			signBit, unsigned))
	}

	if signBitMask := uint32(1) << signBit; (unsigned & signBitMask) == 0 {
		return int32(unsigned)
	}

	mask := -(int32(1) << (signBit + 1))
	return int32(unsigned) | mask
}

func (t immType) parseValue(b InstrBytes) (int32, bool) {
	switch t {
	case immTypeR:
		return 0, false
	case immTypeI:
		val := signExtend(parseBitRange(b, 20, 32), 11)
		return val, true
	case immTypeS:
		low := parseBitRange(b, 7, 12)
		high := parseBitRange(b, 25, 32)
		val := signExtend((high)<<5|low, 11)
		return val, true
	case immTypeB:
		first := parseBitRange(b, 8, 12)
		second := parseBitRange(b, 25, 31)
		third := parseBitRange(b, 7, 8)
		sign := parseBitRange(b, 31, 32)

		// Bit [0] is set to 0 be hardware.
		unsigned := (first << 1) | (second << 5) | (third << 11) | (sign << 12)
		val := signExtend(unsigned, 12)
		return val, true
	case immTypeU:
		unsigned := parseBitRange(b, 12, 32)
		return int32(unsigned << 12), true
	case immTypeJ:
		first := parseBitRange(b, 21, 31)
		second := parseBitRange(b, 20, 21)
		third := parseBitRange(b, 12, 20)
		sign := parseBitRange(b, 31, 32)

		// Bit [0] is set to 0 be hardware.
		unsigned := (first << 1) | (second << 11) | (third << 12) | (sign << 20)
		val := signExtend(unsigned, 20)
		return val, true
	default:
		panic(fmt.Sprintf("unknown immediate type: %v", t))
	}
}
