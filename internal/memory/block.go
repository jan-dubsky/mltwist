package memory

// Block represents a continuous sequence of bytes in a program (user-space)
// address space.
type Block struct {
	Begin Address
	Bytes []byte
}

// NewBlock creates a new block starting at address a and consisting of bytes b.
func NewBlock(a Address, b []byte) Block {
	return Block{
		Begin: a,
		Bytes: b,
	}
}

// Len calculates length of the block in bytes.
func (b Block) Len() int { return len(b.Bytes) }

// End calculates exclusive end of the block.
func (b Block) End() Address { return b.Begin + uint64(b.Len()) }

// Addr returns the longest available slice of memory starting at memory address
// addr.
func (b Block) Addr(addr Address) []byte {
	if addr < b.Begin || addr >= b.End() {
		return nil
	}

	return b.Bytes[addr-b.Begin:]
}
