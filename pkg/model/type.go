package model

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

// typeInvalid represents an invalid value of Type.
const typeInvalid Type = 0

const (
	TypeAritm Type = 1 << iota
	// TypeJump represents an instruction type which execution will
	// unconditionally modify the instruction pointer. In other words, if an
	// instruction i is a jump instruction, it will never happen that
	// execution would follow to instruction offset(i)+len(i), as I will
	// jump elsewhere.
	//
	// Register move instruction is also considered an arithmetic
	// instruction.
	TypeJump
	// TypeCJump represents a conditional jump. In other words a jump which
	// might or might not happen based on a runtime condition.
	TypeCJump
	// TypeJumpDyn represents a jump to dynamic value (value from either
	// register or memory). In other words, it's any jump which target
	// cannot (in generic case) be identified using any sort of static
	// analysis.
	TypeJumpDyn
	// TypeLoad describes an instruction which loads some sort of
	// information from a memory.
	TypeLoad
	// TypeStore describes an instruction which stores some sort of
	// information into a memory.
	TypeStore
	// TypeMemOrder is any instruction which enforces memory ordering. An
	// example of such instructions are for example memory fences.
	TypeMemOrder
	TypeCPUStateChange
	TypeSyscall

	// typeMax is maximal exclusive allowed value of Type. Any value higher
	// or equal to this is invalid.
	typeMax
)

func (t Type) Aritm() bool    { return t&TypeAritm != 0 }
func (t Type) Jump() bool     { return t&TypeJump != 0 }
func (t Type) CJump() bool    { return t&TypeCJump != 0 }
func (t Type) JumpDyn() bool  { return t&TypeJumpDyn != 0 }
func (t Type) Load() bool     { return t&TypeLoad != 0 }
func (t Type) Store() bool    { return t&TypeStore != 0 }
func (t Type) MemOrder() bool { return t&TypeMemOrder != 0 }
func (t Type) Syscall() bool  { return t&TypeSyscall != 0 }
