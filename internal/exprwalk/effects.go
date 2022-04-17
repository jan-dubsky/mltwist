package exprwalk

import (
	"decomp/pkg/expr"
	"fmt"
)

func Effect(effect expr.Effect) []expr.Expr {
	return extractExpr(effect)
}

func Effects(effects ...expr.Effect) []expr.Expr {
	exprs := make([]expr.Expr, 0)
	for _, ef := range effects {
		exprs = append(exprs, extractExpr(ef)...)
	}

	return exprs
}

func extractExpr(effect expr.Effect) []expr.Expr {
	switch ef := effect.(type) {
	case expr.MemStore:
		return []expr.Expr{ef.Addr(), ef.Value()}
	case expr.RegStore:
		return []expr.Expr{ef.Value()}
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", effect))
	}
}
