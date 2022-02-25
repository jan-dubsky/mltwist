package memory

import "decomp/internal/addr"

// Block represents a continuous sequence of bytes in a program (user-space)
// address space.
type Block struct {
	begin addr.Address
	Bytes []byte
}

// NewBlock creates a new block starting at address a and consisting of bytes b.
func NewBlock(a addr.Address, b []byte) Block {
	return Block{
		begin: a,
		Bytes: b,
	}
}

// Begin returns inclusive start address of the block.
func (b Block) Begin() addr.Address { return b.begin }

// Len calculates length of the block in bytes.
func (b Block) Len() int { return len(b.Bytes) }

// End calculates exclusive end of the block.
func (b Block) End() addr.Address { return b.Begin() + uint64(b.Len()) }

// Addr returns the longest available slice of memory starting at memory address
// addr.
func (b Block) Addr(addr addr.Address) []byte {
	if addr < b.Begin() || addr >= b.End() {
		return nil
	}

	return b.Bytes[addr-b.Begin():]
}
