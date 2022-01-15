package instruction

type Instruction struct {
	Type        Type
	ByteLen     uint64
	JumpTargets []Address

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
