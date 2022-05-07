package exprtransform

import (
	"fmt"
	"mltwist/pkg/expr"
)

func JumpAddrs(ex expr.Expr) []expr.Expr {
	jumps := jumpAddrs(ConstFold(ex))
	for i, j := range jumps {
		jumps[i] = ConstFold(j)
	}

	return jumps
}

func jumpAddrs(ex expr.Expr) []expr.Expr {
	switch e := ex.(type) {
	case expr.Binary:
		es1, es2 := jumpAddrs(e.Arg1()), jumpAddrs(e.Arg2())

		es := make([]expr.Expr, 0, len(es1)*len(es2))
		for _, e1 := range es1 {
			for _, e2 := range es2 {
				es = append(es, expr.NewBinary(e.Op(), e1, e2, e.Width()))
			}
		}

		return es
	case expr.Cond:
		es1, es2 := jumpAddrs(e.ExprTrue()), jumpAddrs(e.ExprFalse())

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
		eAddrs := jumpAddrs(e.Addr())

		es := make([]expr.Expr, len(eAddrs))
		for i, eAddr := range eAddrs {
			es[i] = expr.NewMemLoad(e.Key(), eAddr, e.Width())
		}

		return es
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}
