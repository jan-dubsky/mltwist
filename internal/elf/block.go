package elf

import "mltwist/pkg/model"

// Block represents a continuous sequence of bytes in a program (user-space)
// address space.
type Block struct {
	begin      model.Addr
	Bytes      []byte
	executable bool
}

// newBlock creates a new block starting at address a and consisting of bytes b.
func newBlock(a model.Addr, b []byte, executable bool) Block {
	return Block{
		begin:      a,
		Bytes:      b,
		executable: executable,
	}
}

// Begin returns inclusive start address of the block.
func (b Block) Begin() model.Addr { return b.begin }

// Len calculates length of the block in bytes.
func (b Block) Len() int { return len(b.Bytes) }

// End calculates exclusive end of the block.
func (b Block) End() model.Addr { return b.Begin() + model.Addr(b.Len()) }

// Address returns the longest available slice of memory starting at memory
// address a.
func (b Block) Address(a model.Addr) []byte {
	if a < b.Begin() || a >= b.End() {
		return nil
	}

	return b.Bytes[a-b.Begin():]
}

// Executable informs if a block contains executable code.
func (b Block) Executable() bool { return b.executable }
