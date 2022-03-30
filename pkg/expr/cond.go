package expr

var _ Expr = Cond{}

type Cond struct {
	c    Condition
	arg1 Expr
	arg2 Expr

	trueExpr  Expr
	falseExpr Expr

	w Width
}

func NewCond(
	c Condition,
	arg1 Expr,
	arg2 Expr,
	trueExpr Expr,
	falseExpr Expr,
	w Width,
) Cond {
	return Cond{
		c:    c,
		arg1: arg1,
		arg2: arg2,

		trueExpr:  trueExpr,
		falseExpr: falseExpr,

		w: w,
	}
}

func (c Cond) Width() Width { return c.w }
func (Cond) internal()      {}

type Condition uint8

const (
	Lt Condition = iota + 1
	Le
	Eq
	Ne
	Ltu
	Leu
)
