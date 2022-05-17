package model

// Type respresents a special role of instruction in machine code.
//
// Some instruction have a special meaning in in instruction code which cannot
// be represented using an expression model. For those instruction, we introduce
// special markers. One instruction cat fit multiple markers. For this reason,
// we represent a single instruction by a bit mask of its special features.
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

	// TypeMax is maximal exclusive allowed value of Type. Any value higher
	// or equal to this is invalid.
	TypeMax
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
