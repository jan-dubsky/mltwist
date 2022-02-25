package instruction

import "decomp/internal/addr"

type Instruction struct {
	Type    Type
	ByteLen uint64

	// TODO: Consider moving address to another level of abstraction as this
	// is architecture independent value.
	Address addr.Address
	Bytes   []byte

	JumpTargets []addr.Address

	InputMemory   []addr.Address
	InputRegistry []Register

	OutputMemory   []addr.Address
	OutputRegistry []Register

	Details PlatformDetails
}

type PlatformDetails interface {
	// String returns a full string representation of an instruction in
	// assembler code.
	//
	// The representation contains not just the instruction, but also all
	// the registers and memory addresses. All the text should follow
	// platform specific standars how to write instructions and operands.
	String() string
}
