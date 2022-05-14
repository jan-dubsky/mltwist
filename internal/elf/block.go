package elf

import (
	"fmt"
	"mltwist/pkg/model"
)

// Block represents a continuous sequence of bytes in a program (user-space)
// address space.
type Block struct {
	begin model.Addr
	bytes []byte
}

// newBlock creates a new block starting at address a and consisting of bytes b.
func newBlock(a model.Addr, b []byte) Block {
	return Block{
		begin: a,
		bytes: b,
	}
}

// Begin returns inclusive start address of the block.
func (b Block) Begin() model.Addr { return b.begin }

// Len calculates length of the block in bytes.
func (b Block) Len() int { return len(b.bytes) }

// Bytes returns read-only slice of raw bytes stored in the block.
func (b Block) Bytes() []byte { return b.bytes }

// End calculates exclusive end of the block.
func (b Block) End() model.Addr { return b.Begin() + model.Addr(b.Len()) }

// Address returns the longest available slice of memory starting at memory
// address a.
func (b Block) Address(a model.Addr) []byte {
	if a < b.Begin() || a >= b.End() {
		return nil
	}

	return b.bytes[a-b.Begin():]
}

func joinBlocks(b1, b2 Block) Block {
	if e, b := b1.End(), b2.Begin(); e != b {
		panic(fmt.Sprintf("blocks do not follow one another: 0x%x != 0x%x", e, b))
	}

	bytes := make([]byte, len(b1.bytes)+len(b2.bytes))
	copy(bytes, b1.bytes)
	copy(bytes[len(b1.bytes):], b2.bytes)

	return newBlock(b1.Begin(), bytes)
}
