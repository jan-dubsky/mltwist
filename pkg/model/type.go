package model

// Type respresents a role of instruction in machine code.
//
// To implement a generic instruction analysis, we need to be able to say the
// purpose of an instruction in machine code. For this reason, we have to
// introduce some categories of instructions which we will use for this purpose.
//
// As any CISC instruction can be described as a chain of RISC instruction, it
// is in practice sufficient to use RISC categories of instructions. We will
// then use multiple types to describe CISC instruction.
//
// As one instruction might belong into multiple categories, we represent
// instruction type as set of bit flags, where every bit represents a single
// category. This form will allow us to represent an arbitrary set of types for
// every instruction. Zero value of Type then represents no/invalid instruction
// type.
type Type uint64

// TypeNone describes an instruction which is not special in any way. In other
// words, behaviour of such an instruction can be fully described using expr
// package.
const TypeNone Type = 0

const (
	// TypeMemOrder is any instruction which enforces memory ordering. An
	// example of such instructions are for example memory fences.
	TypeMemOrder Type = 1 << iota
	// TypeCPUStateChange describe any CPU state change that might have an
	// impact on execution of further instructions. A typical example of
	// such an instruction would be control registry change which can effect
	// the way the machine code is interpretted.
	TypeCPUStateChange
	// TypeSyscall is any instruction which calls an operating system. This
	// type describes not only syscall instructions, but also any kind of
	// trap which results in operating system interraction with the
	// userspace program.
	TypeSyscall

	// typeMax is maximal exclusive allowed value of Type. Any value higher
	// or equal to this is invalid.
	typeMax
)

// Is checks if type set t contains type other.
//
// As Type type represents a set of instruction types where every single type is
// represented by a single bit which is either set or unset, this check can be
// effectively performed using a single AND operation. Consequence of this Type
// implementation is that this method makes sense only for Type sets with just
// one bit set - i.e. those defined as constants in this package.
func (t Type) Is(other Type) bool { return t&other != 0 }

func (t Type) MemOrder() bool       { return t.Is(TypeMemOrder) }
func (t Type) CPUStateChange() bool { return t.Is(TypeCPUStateChange) }
func (t Type) Syscall() bool        { return t.Is(TypeSyscall) }
