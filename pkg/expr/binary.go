package expr

var _ Expr = Binary{}

// Binary represents any binary operation. This operation is always applied on
// argument expressions of width w and returns as well expression of width w.
type Binary struct {
	op   BinaryOp
	arg1 Expr
	arg2 Expr
	w    Width
}

// NewBinary creates a new binary operation of type op applied on expressios e1
// and e2 with width w.
func NewBinary(op BinaryOp, e1 Expr, e2 Expr, w Width) Binary {
	return Binary{
		op:   op,
		arg1: e1,
		arg2: e2,
		w:    w,
	}
}

// Op is the operation.
func (b Binary) Op() BinaryOp { return b.op }

// Arg1 returns first operang of a binary operation.
func (b Binary) Arg1() Expr { return b.arg1 }

// Arg2 returns second operang of a binary operation.
func (b Binary) Arg2() Expr { return b.arg2 }

// Width returns width of b.
func (b Binary) Width() Width { return b.w }
func (Binary) internalExpr()  {}

// BinaryOp represent any binary operation the virtual CPU supports.
type BinaryOp uint8

const (
	// Add adds 2 numbers. Due to properties of one's complement negative
	// signed integer encoding, signed and unsigned addition are the same.
	Add BinaryOp = iota + 1
	// Sub subtracts second argument from the first one. As in case of Add,
	// both signed and unsigned subtractions are identical.
	Sub
	// Lsh (logically) shifts first operand left second operand number of
	// bits. Second operand is always understood as unsigned.
	Lsh
	// Rsh (logically) shifts first operand left second operand number of
	// bits. Second operand is always understood as unsigned.
	//
	// Logical shift always inserts zeros to highest positions.
	Rsh

	// Mul implements width bits unsigned multiplication of arguments.
	//
	// Signed multiplication is typically implemented as unsigned
	// multiplication of argument magnitudes and the result is then adjusted
	// based on signs of original arguments. Consequently there is no need
	// to support signed multiplication as it can be implemented using
	// unsigned multiplication and further logical operations this
	// expression set provides.
	Mul
	// Div implements width bits unsigned division of the first argument by
	// the second argument.
	//
	// Please note that signed division can be implemented using unsigned
	// division followed by sign resolution logic.
	//
	// Division by zero doesn't cause any error, but produces result of
	// width w with all bits set. In other words, the resulting unsigned
	// value is maximal possible value of width w.
	//
	// This design of division without zero division exception is inspired
	// by RISC-V design, and the reasoning behind is also very similar to
	// the reasoning in RISC-V manual: As we represent arithmetic operations
	// in a functional way we'd have quite hard times to handle any sort of
	// exceptional behaviour. By well-defining the behaviour in case of
	// division by zero, we can significantly simplify the expression
	// analysis logic.
	//
	// Another reason not to produce a division by zero exception is that
	// the design chosen in in a way the more generic one. On CPU
	// architectures where division by zero should cause an exception, this
	// behaviour can be simply implemented by conditional trap execution. On
	// the other hand we'd have quite hard time to represent behaviour of
	// architectures like RISC-V, if we used the division by zero exception
	// as then we'd need a model of operating system handling this exception
	// and filling appropriate values in registry.
	Div

	// And represent bit-wise and of arguments.
	And
	// Or represent bit-wise or of arguments.
	Or
	// Xor represent bit-wise xor of arguments.
	Xor
)
