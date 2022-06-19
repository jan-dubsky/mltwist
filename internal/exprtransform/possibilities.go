package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

// Possibilities returns all possible values the expression can result in. In
// other words, this function removes all conditions in an expression subtree
// and it returns cartesian product of all possibilities.
//
// In some contexts, it makes sense to understand all possible results f an
// expression. This expands the number of expressions we have significantly, but
// it also allows to constant fold some of resulting expressions which can be
// useful in some contexts. A typical example of use-case for this is analysis
// of jump targets where we care only of the expressions but not if they will or
// won't happen.
func Possibilities(ex expr.Expr) []expr.Expr {
	switch e := ex.(type) {
	case expr.Binary:
		es1, es2 := Possibilities(e.Arg1()), Possibilities(e.Arg2())

		es := make([]expr.Expr, 0, len(es1)*len(es2))
		for _, e1 := range es1 {
			for _, e2 := range es2 {
				es = append(es, expr.NewBinary(e.Op(), e1, e2, e.Width()))
			}
		}

		return es
	case expr.Less:
		es1, es2 := Possibilities(e.ExprTrue()), Possibilities(e.ExprFalse())

		es := make([]expr.Expr, 0, len(es1)+len(es2))
		for _, e1 := range es1 {
			es = append(es, SetWidth(e1, e.Width()))
		}
		for _, e2 := range es2 {
			es = append(es, SetWidth(e2, e.Width()))
		}

		return es
	case expr.Const, expr.RegLoad:
		return []expr.Expr{ex}
	case expr.MemLoad:
		eAddrs := Possibilities(e.Addr())

		es := make([]expr.Expr, len(eAddrs))
		for i, eAddr := range eAddrs {
			es[i] = expr.NewMemLoad(e.Key(), eAddr, e.Width())
		}

		return es
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}
