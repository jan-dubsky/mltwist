package model

import (
	"fmt"
	"mltwist/pkg/expr"
)

type Instruction struct {
	// Type is a set of instruction type categories the instruction belongs
	// to.
	Type Type
	// ByteLen is length if an instruction opcode in bytes.
	ByteLen Address

	Effects []expr.Effect

	// JumpTargets is list of known targets of a jump type instruction.
	//
	// This array must be empty is instruction type os not TypeJump,
	// TypeCJump or TypeJumpDyn. On the other hand this array can be empty
	// even for jump types if the jump target cannot be identified from an
	// instruction opcode.
	JumpTargets []Address

	InputMemory []Address

	OutputMemory []Address

	// Details provides additional platform specific properties of the
	// instruction.
	Details PlatformDetails
}

// PlatformDetails is a platform-specific type providing additional information
// about an Instruction.
type PlatformDetails interface {
	// Name returns name of an instruction in assembler code.
	Name() string

	// String returns a full string representation of an instruction in
	// assembler code.
	//
	// Compared to Name method, the representation doesn't contain just an
	// instruction name, but also all instruction operands. Text returned by
	// this method should follow platform-specific notation of instructions
	// operands used in the specific platform.
	String() string
}

// Validate assert that an Instruction description is valid (makes sense). If
// it's not, this method provides a human readable error describing the problem.
func (i *Instruction) Validate() error {
	if t := i.Type; t == TypeInvalid || t >= typeMax {
		return fmt.Errorf("invalid value of type: 0x%x (%d)", t, t)
	}
	if i.ByteLen == 0 {
		return fmt.Errorf("zero ByteLen makes no sense for an instruction")
	}

	for i, effect := range i.Effects {
		if effect == nil {
			return fmt.Errorf("nil effect at position %d", i)
		}
	}

	isJump := i.Type.Jump() || i.Type.CJump() || i.Type.JumpDyn()
	if !isJump && len(i.JumpTargets) > 0 {
		return fmt.Errorf("jump targets list must be empty for non-jump instructions")
	}

	if i.Details == nil {
		return fmt.Errorf("platform details not set")
	}

	return nil
}
