package exprtransform

import (
	"fmt"
	"mltwist/internal/exprtransform/internal/expreval"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"unsafe"
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
			es = append(es, setWidth(e1, e.Width()))
		}
		for _, e2 := range es2 {
			es = append(es, setWidth(e2, e.Width()))
		}

		return es
	case expr.Const, expr.MemLoad, expr.RegLoad:
		return []expr.Expr{ex}
	default:
		panic(fmt.Sprintf("unknown expr.Expr type: %T", ex))
	}
}

func filterJumpAddr(addr model.Address, jumps []expr.Expr) []expr.Expr {
	w := expr.Width(unsafe.Sizeof(addr))
	skipAddr := expreval.ParseConst(expr.NewConstUint(addr, w))

	filtered := make([]expr.Expr, 0, len(jumps))
	for _, j := range jumps {
		if c, ok := j.(expr.Const); ok {
			if expreval.Eq(expreval.ParseConst(c), skipAddr, w) {
				continue
			}
		}

		filtered = append(filtered, j)
	}

	return filtered
}
