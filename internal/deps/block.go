package deps

import (
	"fmt"
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/pkg/model"
	"sort"
)

// block represents a basic block of instructions.
type block struct {
	begin model.Addr
	end   model.Addr
	seq   []*instruction

	// idx is zero-based index of the block in the program.
	idx int
}

// newBlock parses a non-empty sequence of instructions sorted by their
// in-memory addresses into a block and analyzes dependencies in between
// instructions.
func newBlock(idx int, bbSeq []basicblock.Instruction) *block {
	var length model.Addr
	seq := make([]*instruction, len(bbSeq))
	for i, ins := range bbSeq {
		length += ins.Len()
		seq[i] = newInstruction(ins, i)
	}

	processTrueDeps(seq)
	processAntiDeps(seq)
	processOutputDeps(seq)
	processControlDeps(seq)
	processSpecialDeps(seq)

	return &block{
		begin: bbSeq[0].Addr,
		end:   bbSeq[0].Addr + length,
		seq:   seq,

		idx: idx,
	}
}

// Begin returns starting in-memory address of the block. The address relates to
// the original address space of a binary.
func (b *block) Begin() model.Addr { return b.begin }

// End returns in-memory address of the first byte behind the block. The address
// relates to the original address space of a binary.
func (b *block) End() model.Addr { return b.end }

// Bytes returns number of bytes of all instructions in the block.
func (b *block) Bytes() model.Addr { return b.end - b.begin }

// Len returns number of instructions in b.
func (b *block) Len() int { return len(b.seq) }

// Idx returns index of an instruction in list of basic blocks.
func (b *block) Idx() int { return b.idx }

func (b *block) index(i int) *instruction { return b.seq[i] }

// Move moves instruction in the block from index from to index to. All
// instructions in between from and to are shifted one instruction back or
// forward respectively. This method will fail in case the move violates any
// instruction dependency constraints or if either from or to are not valid
// indices of an instruction in the block.
func (b *block) Move(from int, to int) error {
	if err := b.checkMove(from, to); err != nil {
		return fmt.Errorf("cannot move %d to %d: %w", from, to, err)
	}

	move(b.seq, from, to)
	return nil
}

// checkMove asserts of move of instruction on index from to index to is valid
// move in the block.
func (b *block) checkMove(from int, to int) error {
	if err := checkFromToIndex(from, to, len(b.seq)); err != nil {
		return err
	}

	if from < to {
		if u := b.upperBound(from); u < to {
			return fmt.Errorf("upper bound for move is: %d", u)
		}
	} else if from > to {
		if l := b.lowerBound(from); l > to {
			return fmt.Errorf("lower bound for move is: %d", l)
		}
	}

	return nil
}

// findBound finds an instruction boundary (smallest or greatest instruction
// index) in multiple sets of instructions. The cmpF is a comparison predicate
// used to evaluate if the new value of index is "better" than the current (so
// far found) value.
func findBound(cmpF func(first int, second int) bool, set insSet) int {
	var curr int = -1
	for ins := range set {
		if curr < 0 || cmpF(ins.blockIdx, curr) {
			curr = ins.blockIdx
		}
	}

	return curr
}

// lowerBound finds the lowest possible value of index where i can be moved. If
// there is no such lower bound (i.e. i doesn't depend on any previous
// instruction), this method returns zero index.
func (b *block) lowerBound(i int) int {
	ins := b.index(i)
	idx := findBound(func(i, j int) bool { return i > j }, ins.depsBack)

	if idx < 0 {
		return 0
	}
	return idx + 1
}

// upperBound finds the highest possible of index where i can be moved. If there
// is no such upper bound (i.e. i doesn't depend on any later instruction), this
// method returns b.Len() - 1.
func (b *block) upperBound(i int) int {
	ins := b.index(i)
	idx := findBound(func(i, j int) bool { return i < j }, ins.depsFwd)

	if idx < 0 {
		return b.Len() - 1
	}
	return idx - 1
}

func (b *block) setIndex(i int) { b.idx = i }

// setAddr is an empty implementation of address setter which allows us to use
// the same algorithm for both instructions and blocks. As blocks are not
// allowed to move in memory. no other (than empty) implementation of this
// function makes sense.
func (b *block) setAddr(_ model.Addr) {}

// Addr returns the starting address of a basic block.
func (b *block) Addr() model.Addr { return b.begin }

// NextAddr returns an address following the the block.
func (b *block) NextAddr() model.Addr { return b.end }

func (b *block) Address(a model.Addr) (Instruction, bool) {
	i := sort.Search(len(b.seq), func(i int) bool {
		return b.seq[i].Addr() >= a
	})
	if i == len(b.seq) {
		return Instruction{}, false
	}

	ins := b.seq[i]
	if ins.Addr() != a {
		return Instruction{}, false
	}

	return wrapInstruction(ins), true
}
