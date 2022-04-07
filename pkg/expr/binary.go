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

// Width returns width of this operation.
func (b Binary) Width() Width { return b.w }
func (Binary) internal()      {}

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
	//
	// TODO: Rethink right extension of signed and unsigned values.
	Lsh
	// Rsh (logically) shifts first operand left second operand number of
	// bits. Second operand is always understood as unsigned.
	//
	// Logical shift always inserts zeros to highest positions.
	Rsh
	// RshA (arithmetically) shifts first operand left second operand number
	// of bits. Second operand is always understood as unsigned.
	//
	// Arithmetical shift adds copies if the original highest bit to high
	// bit positions.
	RshA

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
	// Signed division can be implemented using unsigned division followed
	// by sign resolution logic.
	Div
	// Mod returns width bits unsigned division reminder of the first
	// argument divided by the second argument.
	//
	// Signed module can be implemented using unsigned module followed by
	// sign resolution logic.
	Mod

	// And represent bit-wise and of arguments.
	And
	// Or represent bit-wise or of arguments.
	Or
	// Xor represent bit-wise xor of arguments.
	Xor
)
