package emulator

import (
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

type Instruction interface {
	Effects() []expr.Effect
	Len() model.Addr
}

// StateProvider provides any state of the program which is not known to a
// program emulation.
type StateProvider interface {
	// Register returns a value of register key of width w.
	//
	// Please note that the expr.Width parameter is necessary as the value
	// returned might be a negative integer. In such a case, the value
	// returned must have exactly width w.
	Register(key expr.Key, w expr.Width) expr.Const

	// Memory returns a value stored in memory identified by key at address
	// addr. The width of the value returned is w.
	Memory(key expr.Key, addr model.Addr, w expr.Width) expr.Const
}
