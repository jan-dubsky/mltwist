package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
	"fmt"
)

type Block struct {
	begin model.Address
	end   model.Address
	seq   []*instruction

	idx int
}

// newBlock parses a non-empty sequence of instructions sorted by their
// in-memory addresses into a Block and analyzes dependencies in between
// instructions.
func newBlock(idx int, seq []repr.Instruction) *Block {
	var length model.Address
	instrs := make([]*instruction, len(seq))
	for i, ins := range seq {
		length += ins.ByteLen
		instrs[i] = newInstruction(ins, i)
	}

	processTrueDeps(instrs)
	processAntiDeps(instrs)
	processOutputDeps(instrs)

	return &Block{
		begin: seq[0].Address,
		end:   seq[0].Address + length,
		seq:   instrs,
		idx:   idx,
	}
}

// Begin returns starting in-memory address of the block. The address relates to
// the original address space of a binary.
func (b *Block) Begin() model.Address { return b.begin }

// End returns in-memory address of the first byte behind the block. The address
// relates to the original address space of a binary.
func (b *Block) End() model.Address { return b.end }

// Bytes returns number of bytes of all instructions in the block.
func (b *Block) Bytes() model.Address { return b.end - b.begin }

// Len returns number of instructions in b.
func (b *Block) Len() int { return len(b.seq) }

// Idx returns index of an instruction in list of basic blocks.
func (b *Block) Idx() int { return b.idx }

// Instructions lists all instructions in b.
func (b *Block) Instructions() []Instruction {
	seq := make([]Instruction, len(b.seq))
	for i, ins := range b.seq {
		seq[i] = ins.ptr()
	}
	return seq
}

// Index returns instruction at index i in b.
func (b *Block) Index(i int) Instruction { return b.seq[i].ptr() }

// Move moves instruction in the block from index from to index to. All
// instructions in between from and to are shifted one instruction back or
// forward respectively. This method will fail in case the move violates
// instruction dependency constraints.
func (b *Block) Move(from int, to int) error {
	if err := b.checkMove(from, to); err != nil {
		return fmt.Errorf("cannot move %d to %d: %w", from, to, err)
	}

	if from == to {
		return nil
	}

	f := b.seq[from]
	if from < to {
		b.moveFwd(from, to)
	} else {
		b.moveBack(from, to)
	}
	b.seq[to] = f
	f.blockIdx = to

	return nil
}

func (b *Block) moveFwd(from int, to int) {
	for i := from; i < to; i++ {
		b.seq[i] = b.seq[i+1]
		b.seq[i].blockIdx = i
	}
}

func (b *Block) moveBack(from int, to int) {
	for i := from; i > to; i-- {
		b.seq[i] = b.seq[i-1]
		b.seq[i].blockIdx = i
	}
}

func (b *Block) validateIndex(name string, value int) error {
	if value < 0 {
		return fmt.Errorf("negative value of %q is not allowed: %d", name, value)
	}
	if l := len(b.seq); value >= l {
		return fmt.Errorf("value of %q is above limit: %d >= %d", name, value, l)
	}

	return nil
}

func (b *Block) checkMove(from int, to int) error {
	if err := b.validateIndex("from", from); err != nil {
		return err
	} else if err := b.validateIndex("to", to); err != nil {
		return err
	}

	if from < to {
		if u := b.upperBound(b.seq[from]); u < to {
			return fmt.Errorf("upper bound for move is: %d", u)
		}
	} else if from > to {
		if l := b.lowerBound(b.seq[from]); l > to {
			return fmt.Errorf("lower bound for move is: %d", l)
		}
	}

	return nil
}

func findOne(f func(int, int) bool, instrs insSet, idx int) int {
	curr := idx
	for ins := range instrs {
		if curr < 0 || f(ins.blockIdx, curr) {
			curr = ins.blockIdx
		}
	}

	return curr
}

// findBound finds an instruction boundary (smallest or greatest instruction
// index) in multiple sets of instructions. The f predicate is in the search
// always used to evaluate if the new value of index is "better" than the
// current (so far found) value.
func findBound(f func(int, int) bool, sets ...insSet) int {
	var curr int = -1
	for _, s := range sets {
		curr = findOne(f, s, curr)
	}

	return curr
}

// LowerBound finds the lowest possible value of index where i can be moved. If
// there is no such lower bound (i.e. i doesn't depend on any previous
// instruction), this method returns zero index.
func (b *Block) LowerBound(i Instruction) int { return b.lowerBound(i.i) }

func (*Block) lowerBound(i *instruction) int {
	idx := findBound(func(i, j int) bool { return i > j },
		i.trueDepsBack,
		i.antiDepsBack,
		i.outputDepsBack,
	)

	if idx < 0 {
		return 0
	}
	return idx + 1
}

// UpperBound finds the highest possible of index where i can be moved. If there
// is no such upper bound (i.e. i doesn't depend on any later instruction), this
// method returns b.Len() - 1.
func (b *Block) UpperBound(i Instruction) int { return b.upperBound(i.i) }

func (b *Block) upperBound(i *instruction) int {
	idx := findBound(func(i, j int) bool { return i < j },
		i.trueDepsFwd,
		i.antiDepsFwd,
		i.outputDepsFwd,
	)

	if idx < 0 {
		return b.Len() - 1
	}
	return idx - 1
}
