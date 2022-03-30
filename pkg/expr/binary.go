package expr

var _ Expr = Binary{}

type Binary struct {
	op   BinaryOp
	arg1 Expr
	arg2 Expr
	w    Width
}

func NewBinary(op BinaryOp, e1 Expr, e2 Expr, w Width) Binary {
	return Binary{
		op:   op,
		arg1: e1,
		arg2: e2,
		w:    w,
	}
}

func (b Binary) Width() Width { return b.w }
func (Binary) internal()      {}

type BinaryOp uint8

const (
	Add BinaryOp = iota + 1
	Sub
	Mul
	Div
	Mod
	Rsh
	Lsh
	RshL
	And
	Or
	Xor
)
