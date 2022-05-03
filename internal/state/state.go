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
		s.Regs[e.Key()] = exprtransform.SetWidth(e.Value(), e.Width())
		return true
	default:
		panic(fmt.Sprintf("unknown expr.Effect type: %T", ef))
	}
}
