package instruction

// Type respresents a role of instruction in any kind of machine code.
//
// To implement any sort of generic instruction analysis, we need to be able to
// say the purpose of an instruction in machine code. For this reason, we have
// to introduce some categories of instructions which we will use for this
// purpose.
//
// It seems to be a good idea to use RISC-like categories as most of CISC
// instructions typically correctpond to some chain of RISC instructions.
//
// As one CISC instruction might belong into multiple categories, we represent
// instruction type as set of bit flags, where every bit represents a single
// category. This form will allow us to represent an arbitrary set of categories
// for every instruction. Zero value of Type then represents no instruction
// type.
type Type uint64

const (
	TypeMov Type = 1 << iota
	TypeAritm
	// TypeJump represents an instruction type which execution will
	// unconditionally modify the instruction pointer. In other words, if an
	// instruction i is a jump instruction, it will never happen that
	// execution would follow to instruction offset(i)+len(i), as I will
	// jump elsewhere.
	TypeJump
	TypeCJump
	TypeLoad
	TypeStore
	TypeMemOrder
	TypeCPUStateChange
	TypeSyscall
)

func (t Type) Mov() bool   { return t&TypeMov != 0 }
func (t Type) Aritm() bool { return t&TypeAritm != 0 }
func (t Type) Jump() bool  { return t&TypeJump != 0 }
func (t Type) CJump() bool { return t&TypeCJump != 0 }
func (t Type) Load() bool  { return t&TypeLoad != 0 }
func (t Type) Store() bool { return t&TypeStore != 0 }
