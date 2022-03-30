package expr

var _ Expr = Unary{}

type Unary struct {
	op  UnaryOp
	arg Expr
	w   Width
}

func NewUnary(op UnaryOp, e Expr, w Width) Unary {
	return Unary{
		op:  op,
		arg: e,
		w:   w,
	}
}

func (u Unary) Width() Width { return u.w }
func (Unary) internal()      {}

type UnaryOp uint8

const (
	Negate UnaryOp = iota + 1
	Not
)
