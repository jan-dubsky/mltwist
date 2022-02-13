package riscv

// addOpcodeInfo describes any sort of additional information stored directly in
// an instruction opcode, which is not an opcode itself, neither an immediate
// value.
//
// Even though RISC-V specification itself states that there are only opcodes
// and 5 types of immediate value encoding, there are in reality instructions
// which containt additional information in the opcode. The RISC-V specification
// itself doesn't name any specific mechanism how those additional values are
// handled hor how they should be reasoned, it's useful for us to group them
// under one concept and that is additional opcode information concept.
type addOpcodeInfo uint8

const (
	// addImmSh32 is an additional bit shift information stored in an
	// instruction opcode. The data format itself is a modification of
	// immTypeI used to encode bit shift instructions with bit shift encoded
	// in an immediate value of the opcode for 32 bit registers.
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
	addImmSh32 addOpcodeInfo = iota + 1

	// addImmSh64 is an additional bit shift information stored in an
	// instruction opcode. The data format itself is a modification of
	// immTypeI used to encode bit shift instructions with bit shift encoded
	// in an immediate value of the opcode for 64 bit registers.
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
	addImmSh64

	// addImmCSR is an additional CCR register number stored directly in an
	// instruction opcode.
	//
	// This format of data encoding is very similar to immTypeI, where CSR
	// number is stored exactly where an immediate value is stored in an
	// opcode. The only reason why we don't consider CSR to be exactly an
	// immediate value is the fact that there exist instructions which write
	// an immediate value into a CSR. So if we considered CSR number to be
	// an immediate value, we'd get into a situation where we have 2 distinc
	// immediate values. As such a setup would be super confusing, we rather
	// decided to treat CSR as an exceptional (additional opcode
	// information) to eliminate the irregularity asm uch as possible.
	addImmCSR
)
