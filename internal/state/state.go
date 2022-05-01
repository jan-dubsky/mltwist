package state

import (
	"fmt"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type State struct {
	// Regs represents state of registry file.
	Regs RegMap
	// Mems is set of memories and their respective states.
	Mems MemMap
}

// New creates a new state.
func New() *State {
	return &State{
		// We need to set some value as 100 registers (default map size
		// in Go) is too many. So 32 is thumbsucked, but more reasonable
		// value.
		Regs: make(RegMap, 32),
		// We assume to have 1 address space - the real program memory.
		Mems: make(MemMap, 1),
	}
}

// Apply changes state by applying effect ef.
func (s *State) Apply(ef expr.Effect) {
	switch e := ef.(type) {
	case expr.MemStore:
		c, ok := e.Addr().(expr.Const)
		if !ok {
			return
		}

		addr, _ := expr.ConstUint[model.Addr](c)
		s.Mems.Store(e.Key(), addr, e.Value(), e.Width())
	case expr.RegStore:
		s.Regs[e.Key()] = exprtransform.SetWidth(e.Value(), e.Width())
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", ef))
	}
}

// Substitute replaces all register loads and memory loads in ex and all its
// subexpressions by values in Regs and Mems respectively. Constant fonding is
// applied to the expression returned.
//
// If value for a register or memory loaded is in the state, those loads are
// left unchanged. If memory address to load is not a constant expression (or
// cannot be made constant by constant fonding), the memory load will be left
// unchanged as well.
func (s *State) Substitute(ex expr.Expr) expr.Expr {
	ex = exprtransform.ReplaceAll(ex, func(curr expr.RegLoad) (expr.Expr, bool) {
		e, ok := s.Regs[curr.Key()]
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
		ex, ok := s.Mems.Load(curr.Key(), addr, curr.Width())
		if !ok {
			return curr, false
		}

		return ex, true
	})

	return exprtransform.ConstFold(ex)
}
