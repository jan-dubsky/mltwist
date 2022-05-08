package memory

import (
	"mltwist/internal/state/interval"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// MamMap is set of individual memory address spaces identified by a key.
type MemMap map[expr.Key]Memory

// Load loads w bytes from an address addr in address space identified by key.
func (m MemMap) Load(key expr.Key, addr model.Addr, w expr.Width) (expr.Expr, bool) {
	mem, ok := m[key]
	if !ok {
		return nil, false
	}

	return mem.Load(addr, w)
}

// Store stores w bytes of ex to address addr in address space identified by
// key. If address space for key key doesn't exist yet, it's created.
func (m MemMap) Store(key expr.Key, addr model.Addr, ex expr.Expr, w expr.Width) {
	mem, ok := m[key]
	if !ok {
		mem = NewSparse()
		m[key] = mem
	}

	mem.Store(addr, ex, w)
}

func (m MemMap) Missing(key expr.Key, addr model.Addr, w expr.Width) interval.Map[model.Addr] {
	mem, ok := m[key]
	if !ok {
		return interval.NewMap(interval.New(addr, addr+model.Addr(w)))
	}

	return mem.Missing(addr, w)
}

func (m MemMap) Blocks(key expr.Key) interval.Map[model.Addr] {
	mem, ok := m[key]
	if !ok {
		return interval.NewMap[model.Addr]()
	}

	return mem.Blocks()
}
