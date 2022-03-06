package model

import "fmt"

type Instruction struct {
	Type    Type
	ByteLen uint64

	JumpTargets []Address

	InputMemory   []Address
	InputRegistry map[Register]struct{}

	OutputMemory   []Address
	OutputRegistry map[Register]struct{}

	Details PlatformDetails
}

type PlatformDetails interface {
	// String returns a full string representation of an instruction in
	// assembler code.
	//
	// The representation contains not just the instruction, but also all
	// the registers and memory addresses. All the text should follow
	// platform specific notation of instructions operands and immediate.
	String() string
}

// Validate assert that an Instruction description is valid (makes sense). If
// it's not, this method provides a human readable error describing the problem.
func (i *Instruction) Validate() error {
	if t := i.Type; t == typeInvalid || t >= typeMax {
		return fmt.Errorf("invalid value of type: 0x%x (%d)", t, t)
	}

	if i.ByteLen == 0 {
		return fmt.Errorf("zero ByteLen makes no sense for an instruction")
	}

	return nil
}
