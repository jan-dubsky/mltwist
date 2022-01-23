package riscv

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

	// immTypeISh32 is an instruction opcode preudo-format which we defined
	// ourselves - it's not listed in the RISC V specification). This format
	// is a modification of immTypeI used to encode bit shift instructions
	// with bit shift encoded in an immediate value of the opcode for 32 bit
	// registers.
	//
	// This format is standard immTypeI, but as bit shift of more than
	// number of bits in a register width does not make sense, the immediate
	// value is limited to as many bits to express XLEN (5 for 32 bit XLEN).
	// All other bits of immediate value (but one - see below) are reserved.
	//
	// Moreover, the encoding of logical and arithmetic shift opcodes differ
	// just by bit [30] of the opcode, which in I type instruction encoding
	// represents 10th bit of an immediate value. This style of encoding is
	// very similar to a way how R type opcodes encode func7 part of opcode,
	// but it differs because there is no immediate value in R type opcode
	// format. In this pseudo-format, there is still an immediate value
	// encoded, there are just not all 12 bits of immediate used for an
	// immediate. Instead, some of them are reserved and some of them are
	// used to encode a part of instruction opcode.
	immTypeISh32

	// immTypeISh64 is an instruction opcode preudo-format which we defined
	// ourselves - it's not listed in the RISC V specification). This format
	// is a modification of immTypeI used to encode bit shift instructions
	// with bit shift encoded in an immediate value of the opcode for 64 bit
	// registers.
	//
	// This format is standard immTypeI, but as bit shift of more than
	// number of bits in a register width does not make sense, the immediate
	// value is limited to as many bits to express XLEN (6 for 64 bit XLEN).
	// All other bits of immediate value (but one - see below) are reserved.
	//
	// Moreover, the encoding of logical and arithmetic shift opcodes differ
	// just by bit [30] of the opcode, which in I type instruction encoding
	// represents 10th bit of an immediate value. This style of encoding is
	// very similar to a way how R type opcodes encode func7 part of opcode,
	// but it differs because there is no immediate value in R type opcode
	// format. In this pseudo-format, there is still an immediate value
	// encoded, there are just not all 12 bits of immediate used for an
	// immediate. Instead, some of them are reserved and some of them are
	// used to encode a part of instruction opcode.
	immTypeISh64
)

func parseImmediate(tp immType, b []byte) {
	assertOpcodeLen(b)

	// FIXME:
}
