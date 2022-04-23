package memory

import "mltwist/pkg/model"

// Block represents a continuous sequence of bytes in a program (user-space)
// address space.
type Block struct {
	begin model.Addr
	Bytes []byte
}

// NewBlock creates a new block starting at address a and consisting of bytes b.
func NewBlock(a model.Addr, b []byte) Block {
	return Block{
		begin: a,
		Bytes: b,
	}
}

// Begin returns inclusive start address of the block.
func (b Block) Begin() model.Addr { return b.begin }

// Len calculates length of the block in bytes.
func (b Block) Len() int { return len(b.Bytes) }

// End calculates exclusive end of the block.
func (b Block) End() model.Addr { return b.Begin() + model.Addr(b.Len()) }

// Addr returns the longest available slice of memory starting at memory address
// addr.
func (b Block) Addr(addr model.Addr) []byte {
	if addr < b.Begin() || addr >= b.End() {
		return nil
	}

	return b.Bytes[addr-b.Begin():]
}
