package memory

import "decomp/pkg/model"

// Block represents a continuous sequence of bytes in a program (user-space)
// address space.
type Block struct {
	begin model.Address
	Bytes []byte
}

// NewBlock creates a new block starting at address a and consisting of bytes b.
func NewBlock(a model.Address, b []byte) Block {
	return Block{
		begin: a,
		Bytes: b,
	}
}

// Begin returns inclusive start address of the block.
func (b Block) Begin() model.Address { return b.begin }

// Len calculates length of the block in bytes.
func (b Block) Len() int { return len(b.Bytes) }

// End calculates exclusive end of the block.
func (b Block) End() model.Address { return b.Begin() + model.Address(b.Len()) }

// Addr returns the longest available slice of memory starting at memory address
// addr.
func (b Block) Addr(addr model.Address) []byte {
	if addr < b.Begin() || addr >= b.End() {
		return nil
	}

	return b.Bytes[addr-b.Begin():]
}
