package expr

var _ Expr = Less{}

// Less represents a conditional expression comparing arg1 and arg2 using
// unsigned less operation. This operation returns trueExpr or falseExpr if
// condition result is true or false respectively. The width of expression
// returned is always w.
//
// Please note that any other comparison as signed comparison of equality
// comparison can be achieved by using multiple unsigned Less integer
// comparisons.
type Less struct {
	arg1 Expr
	arg2 Expr

	trueExpr  Expr
	falseExpr Expr

	w Width
}

// NewLess returns new Less matching comparing arg1 and arg2 and returning
// trueExpr or falseExpr for true of false unsigned integer comparison result
// respectively. The width of expression returned is always w.
func NewLess(
	arg1 Expr,
	arg2 Expr,
	trueExpr Expr,
	falseExpr Expr,
	w Width,
) Less {
	return Less{
		arg1: arg1,
		arg2: arg2,

		trueExpr:  trueExpr,
		falseExpr: falseExpr,

		w: w,
	}
}

// Arg1 returns first argument of the less comparison.
func (l Less) Arg1() Expr { return l.arg1 }

// Arg2 returns second argument of the less comparison.
func (l Less) Arg2() Expr { return l.arg2 }

// ExprTrue returns the expression returned in case of condition being true.
func (l Less) ExprTrue() Expr { return l.trueExpr }

// ExprTrue returns the expression returned in case of condition being false.
func (l Less) ExprFalse() Expr { return l.falseExpr }

// Width returns width of l.
func (l Less) Width() Width { return l.w }
func (Less) internalExpr()  {}
