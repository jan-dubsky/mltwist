package deps

// Instruction represents a single instruction in the code.
type Instruction struct {
	*instruction
}

func wrapInstruction(ins *instruction) Instruction { return Instruction{ins} }

// String returns string representation of an instruction. This representation
// follows standard platform-specific way of assembler code syntax for a given
// platform.
func (i Instruction) String() string { return i.details.String() }

// Block represents a single basic-block in the code.
type Block struct {
	*block
}

func wrapBlock(b *block) Block { return Block{b} }

// Instructions lists all instructions in b.
func (b Block) Instructions() []Instruction {
	seq := make([]Instruction, len(b.seq))
	for i, ins := range b.seq {
		seq[i] = wrapInstruction(ins)
	}
	return seq
}

// Index returns instruction at index i in b.
func (b Block) Index(i int) Instruction { return wrapInstruction(b.index(i)) }
