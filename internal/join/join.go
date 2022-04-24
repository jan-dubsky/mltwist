package join

import (
	"fmt"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

func Join(ops [][]expr.Effect) [][]expr.Effect {
	j := NewJoiner()

	joined := make([][]expr.Effect, len(ops))
	for i, op := range ops {
		joined[i] = j.Add(op)
	}

	return joined
}

type Joiner struct {
	regs regMap
	mems memMap
}

func NewJoiner() *Joiner {
	return &Joiner{
		// We need to set some value as 100 registers (default map size
		// in Go) is too many. So 32 is thumbsucked, but more reasonable
		// value.
		regs: make(regMap, 32),
		// We assume to have 1 address space - the real program memory.
		mems: make(memMap, 1),
	}
}

func (j *Joiner) Add(effs []expr.Effect) []expr.Effect {
	substituted := make([]expr.Effect, len(effs))
	for i, ef := range effs {
		substituted[i] = j.Substitute(ef)
	}

	for _, s := range substituted {
		j.apply(s)
	}

	return substituted
}

func (j *Joiner) Substitute(ef expr.Effect) expr.Effect {
	switch e := ef.(type) {
	case expr.MemStore:
		value := exprtransform.ConstFold(j.SubstituteExpr(e.Value()))
		addr := exprtransform.ConstFold(j.SubstituteExpr(e.Addr()))
		return expr.NewMemStore(value, e.Key(), addr, e.Width())
	case expr.RegStore:
		value := exprtransform.ConstFold(j.SubstituteExpr(e.Value()))
		return expr.NewRegStore(value, e.Key(), e.Width())
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", ef))
	}
}

func (j *Joiner) SubstituteExpr(ex expr.Expr) expr.Expr {
	ex = exprtransform.ReplaceAll(ex, func(curr expr.RegLoad) (expr.Expr, bool) {
		e, ok := j.regs[curr.Key()]
		if !ok {
			return curr, false
		}

		return exprtransform.SetWidth(e, curr.Width()), true
	})

	ex = exprtransform.ReplaceAll(ex, func(curr expr.MemLoad) (expr.Expr, bool) {
		c, ok := exprtransform.ConstFold(curr.Addr()).(expr.Const)
		if !ok {
			return curr, false
		}

		addr, _ := expr.ConstUint[model.Addr](c)
		ex, ok := j.mems.load(curr.Key(), addr, curr.Width())
		if !ok {
			return curr, false
		}

		return ex, true
	})

	return ex
}

func (j *Joiner) apply(ef expr.Effect) {
	switch e := ef.(type) {
	case expr.MemStore:
		c, ok := e.Addr().(expr.Const)
		if !ok {
			return
		}

		addr, _ := expr.ConstUint[model.Addr](c)
		j.mems.store(e.Key(), addr, e.Value(), e.Width())
	case expr.RegStore:
		j.regs[e.Key()] = exprtransform.SetWidth(e.Value(), e.Width())
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", ef))
	}
}
