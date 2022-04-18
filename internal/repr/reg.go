package repr

import (
	"mltwist/internal/exprtransform"
	"mltwist/internal/exprwalk"
	"mltwist/pkg/expr"
)

type RegSet map[string]struct{}

type inputRegCollector struct {
	exprwalk.EmptyExprWalker
	inputRegs RegSet
}

func (c *inputRegCollector) RegLoad(e expr.RegLoad) error {
	c.inputRegs[string(e.Key())] = struct{}{}
	return nil
}

func regs(effects []expr.Effect) (RegSet, RegSet) {
	c := &inputRegCollector{inputRegs: make(RegSet, 2)}

	for _, ex := range exprwalk.Effects(effects...) {
		_ = exprwalk.Expr(c, ex)
	}

	outputRegs := make(RegSet, 1)
	for _, ef := range effects {
		store, ok := ef.(expr.RegStore)
		if !ok {
			continue
		}

		outputRegs[string(store.Key())] = struct{}{}
	}

	return c.inputRegs, outputRegs
}

func jumps(effects []expr.Effect) []expr.Expr {
	var jumpAddrs []expr.Expr
	for _, ef := range effects {
		e, ok := ef.(expr.RegStore)
		if !ok {
			continue
		}

		if e.Key() != expr.IPKey {
			continue
		}

		addrs := exprtransform.JumpAddrs(e.Value())
		jumpAddrs = append(jumpAddrs, addrs...)
	}

	return jumpAddrs
}
