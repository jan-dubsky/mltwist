package expr

type Expr interface{}

type Unary struct {
	op  UnaryOp
	arg Expr
}

func NewUnary(op UnaryOp, e Expr) Unary {
	return Unary{
		op:  op,
		arg: e,
	}
}

func (u Unary) Op() UnaryOp { return u.op }

type Binary struct {
	op   BinaryOp
	arg1 Expr
	arg2 Expr
}

func NewBinary(op BinaryOp, e1 Expr, e2 Expr) Binary {
	return Binary{
		op:   op,
		arg1: e1,
		arg2: e2,
	}
}
