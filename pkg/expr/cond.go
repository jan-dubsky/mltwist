package expr

var _ Expr = Cond{}

// Cond represents a conditional expression comparing arg1 and arg2. This
// operation returns trueExpr or falseExpr if condition result is true or false
// respectively. The width of expression returned is always w.
type Cond struct {
	cond Condition
	arg1 Expr
	arg2 Expr

	trueExpr  Expr
	falseExpr Expr

	w Width
}

// NewCond returns new Cond matching cond for arg1 and arg2 and returning
// trueExpr or falseExpr for true of false condition result respectively. The
// width of expression returned is always w.
func NewCond(
	cond Condition,
	arg1 Expr,
	arg2 Expr,
	trueExpr Expr,
	falseExpr Expr,
	w Width,
) Cond {
	return Cond{
		cond: cond,
		arg1: arg1,
		arg2: arg2,

		trueExpr:  trueExpr,
		falseExpr: falseExpr,

		w: w,
	}
}

// Cond returns condidion applies on Arg1 and Arg2.
func (c Cond) Cond() Condition { return c.cond }

// Arg1 returns first argument of a condition.
func (c Cond) Arg1() Expr { return c.arg1 }

// Arg2 returns second argument of a condition.
func (c Cond) Arg2() Expr { return c.arg2 }

// ExprTrue returns the expression returned in case of condition being true.
func (c Cond) ExprTrue() Expr { return c.trueExpr }

// ExprTrue returns the expression returned in case of condition being false.
func (c Cond) ExprFalse() Expr { return c.falseExpr }

// Width returns width of c.
func (c Cond) Width() Width { return c.w }
func (Cond) internalExpr()  {}

// Condition represents a comparison applied on Cond expression arguments.
type Condition uint8

const (
	// Eq performes equality (==) comparison of 2 arguments.
	Eq Condition = iota + 1
	// Ltu performs unsigned less-then (<) comparison of the first and
	// second argument.
	Ltu
	// Leu performs unsigned less-then-or-equal (<=) comparison of the first
	// and second argument.
	Leu
	// Lts perform signed less-then (<) comparison of the first and second
	// argument.
	Lts
	// Les performs signed less-then-or-equal (<=) comparison of the first
	// and second argument.
	Les
)
