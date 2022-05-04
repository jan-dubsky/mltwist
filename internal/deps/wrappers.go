package deps

// Instruction represents a single instruction in the program.
type Instruction struct {
	*instruction
}

func wrapInstruction(ins *instruction) Instruction { return Instruction{ins} }

// String returns string representation of an instruction. This representation
// follows standard platform-specific way of assembler code syntax for a given
// platform.
func (i Instruction) String() string { return i.details.String() }

// Block represents a single basic-block in the program.
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

// LowerBound finds the lowest possible value of index where i can be moved. If
// there is no such lower bound (i.e. i doesn't depend on any previous
// instruction), this method returns zero index.
func (b Block) LowerBound(i int) int { return b.lowerBound(i) }

// UpperBound finds the highest possible of index where i can be moved. If there
// is no such upper bound (i.e. i doesn't depend on any later instruction), this
// method returns b.Len() - 1.
func (b Block) UpperBound(i int) int { return b.upperBound(i) }
