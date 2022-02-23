package instruction

import "decomp/internal/memory"

type Instruction struct {
	Type    Type
	ByteLen uint64

	JumpTargets []memory.Address

	InputMemory   []memory.Address
	InputRegistry []Register

	OutputMemory   []memory.Address
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
