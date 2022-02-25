package memory

import (
	"decomp/internal/addr"
	"decomp/internal/interval"
)

// Memory is sparse representation of program (user-space) memory space.
//
// Memory is allowed to represent only partition of the real user-space memory.
// For example, it might represent only certain memory types (stack, static
// data, code). The exact meaning of any memory structure is given by the usage
// and not by the struct type itself.
//
// This type has to be created using New method.
type Memory struct {
	// Blocks is list of individual memory blocks which are sorted in
	// ascending offset order and which are non-overlapping.
	//
	// After Memory object creation, this field has to be treated as
	// immutable.
	Blocks []Block

	l *interval.List
}

// New creates a new memory structure. This method return an error if blocks
// overlap.
func New(blocks ...Block) (*Memory, error) {
	intervals := make([]interval.Interval, len(blocks))
	for i, b := range blocks {
		intervals[i] = b
	}

	l, err := interval.NewList(intervals)
	if err != nil {
		return nil, err
	}

	return &Memory{Blocks: blocks, l: l}, nil
}

// Addr returns the longest available slice of memory starting at memory address
// addr.
func (m *Memory) Addr(addr addr.Address) []byte {
	b := m.l.Addr(addr)
	if b == nil {
		return nil
	}

	return b.(Block).Addr(addr)
}
