package deps

import (
	"decomp/internal/repr"
	"fmt"
)

type Block struct {
	seq []*instruction
}

func newBlock(seq []repr.Instruction) *Block {
	instrs := make([]*instruction, len(seq))
	for i, ins := range seq {
		instr := newInstruction(ins)
		instr.blockIdx = i
		instrs[i] = instr
	}

	processTrueDeps(instrs)
	processAntiDeps(instrs)
	processOutputDeps(instrs)

	return &Block{seq: instrs}
}

// Len returns number of instructions in b.
func (b *Block) Len() int { return len(b.seq) }

// Instructions lists all instructions in b.
func (b *Block) Instructions() []Instruction {
	seq := make([]Instruction, len(b.seq))
	for i, ins := range b.seq {
		seq[i] = ins.ptr()
	}
	return seq
}

// Idx returns instruction at index i in b.
func (b *Block) Idx(i int) Instruction { return b.seq[i].ptr() }

func (b *Block) validateIndex(name string, value int) error {
	if value < 0 {
		return fmt.Errorf("negative value of %q is not allowed: %d", name, value)
	}
	if l := len(b.seq); value >= l {
		return fmt.Errorf("value of %q is above limit: %d >= %d", name, value, l)
	}

	return nil
}

func (b *Block) Move(from int, to int) error {
	if err := b.validateIndex("from", from); err != nil {
		return err
	} else if err := b.validateIndex("to", to); err != nil {
		return err
	} else if from == to {
		return nil
	}

	// TODO: Complete this function.

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
