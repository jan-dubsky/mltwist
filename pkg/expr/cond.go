package expr

type Condition uint8

const (
	Lt Condition = iota + 1
	Le
	Eq
	Ne
	Ge
	Gt
)

type Cond struct {
	c    Condition
	arg1 Expr
	arg2 Expr

	trueExpr  Expr
	falseExpr Expr
}

func NewCond(c Condition, arg1 Expr, arg2 Expr, trueExpr Expr, falseExpr Expr) Cond {
	return Cond{
		c:    c,
		arg1: arg1,
		arg2: arg2,

		trueExpr:  trueExpr,
		falseExpr: falseExpr,
	}
}
