package state

import (
	"fmt"
	"mltwist/internal/exprtransform"
	"mltwist/internal/state/memory"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type State struct {
	// Regs represents state of registry file.
	Regs *RegMap
	// Mems is set of memories and their respective states.
	Mems memory.MemMap
}

// New creates a new state.
func New() *State {
	return &State{
		Regs: NewRegMap(),
		// We assume to have 1 address space - the real program memory.
		Mems: make(memory.MemMap, 1),
	}
}

// Apply changes state by applying effect ef and provides an information if it
// was possible to apply the effect.
func (s *State) Apply(ef expr.Effect) bool {
	switch e := ef.(type) {
	case expr.MemStore:
		c, ok := exprtransform.ConstFold(e.Addr()).(expr.Const)
		if !ok {
			return false
		}

		addr, _ := expr.ConstUint[model.Addr](c)
		s.Mems.Store(e.Key(), addr, e.Value(), e.Width())
		return true
	case expr.RegStore:
		s.Regs.Store(e.Key(), e.Value(), e.Width())
		return true
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", ef))
	}
}
