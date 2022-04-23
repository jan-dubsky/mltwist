package memory

import (
	"fmt"
	"mltwist/pkg/model"
	"sort"
)

// Memory is sparse representation of program (user-space) memory space.
//
// Memory is allowed to represent only partition of the real user-space memory.
// For example, it might represent only certain memory types (stack, static
// data, code). The exact meaning of any memory structure is given by the usage
// and not by the struct type itself.
type Memory struct {
	// Blocks is list of individual memory blocks which are sorted in
	// ascending offset order and which are non-overlapping.
	Blocks []Block
}

// New creates a new memory structure. This method return an error if blocks
// overlap.
func New(blocks ...Block) (*Memory, error) {
	if len(blocks) == 0 {
		return &Memory{}, nil
	}

	less := func(i, j int) bool { return blocks[i].Begin() < blocks[j].Begin() }
	sort.Slice(blocks, less)

	for i := range blocks[1:] {
		if e, b := blocks[i].End(), blocks[i+1].Begin(); e > b {
			return nil, fmt.Errorf(
				"block %d (ending 0x%x) and %d (starting 0x%x) overlap",
				i, e, i+1, b)
		}
	}

	return &Memory{Blocks: blocks}, nil
}

// Addr returns the longest available slice of memory starting at memory address
// addr.
func (m *Memory) Addr(addr model.Addr) []byte {
	idx := sort.Search(len(m.Blocks), func(i int) bool {
		return m.Blocks[i].End() > addr
	})

	if idx == len(m.Blocks) || m.Blocks[idx].Begin() > addr {
		return nil
	}

	return m.Blocks[idx].Addr(addr)
}
