package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
	"fmt"
)

type block struct {
	begin model.Address
	end   model.Address
	seq   []*instruction

	// idx is zero-based index of the block in the program.
	idx int
}

// newBlock parses a non-empty sequence of instructions sorted by their
// in-memory addresses into a block and analyzes dependencies in between
// instructions.
func newBlock(idx int, seq []repr.Instruction) *block {
	var length model.Address
	instrs := make([]*instruction, len(seq))
	for i, ins := range seq {
		length += ins.ByteLen
		instrs[i] = newInstruction(ins, i)
	}

	processTrueDeps(instrs)
	processAntiDeps(instrs)
	processOutputDeps(instrs)
	processControlDeps(instrs)
	processSpecialDeps(instrs)

	return &block{
		begin: seq[0].Address,
		end:   seq[0].Address + length,
		seq:   instrs,
		idx:   idx,
	}
}

// Begin returns starting in-memory address of the block. The address relates to
// the original address space of a binary.
func (b *block) Begin() model.Address { return b.begin }

// End returns in-memory address of the first byte behind the block. The address
// relates to the original address space of a binary.
func (b *block) End() model.Address { return b.end }

// Bytes returns number of bytes of all instructions in the block.
func (b *block) Bytes() model.Address { return b.end - b.begin }

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
func findBound(cmpF func(first int, second int) bool, sets ...insSet) int {
	var curr int = -1
	for _, s := range sets {
		for ins := range s {
			if curr < 0 || cmpF(ins.blockIdx, curr) {
				curr = ins.blockIdx
			}
		}
	}

	return curr
}

// lowerBound finds the lowest possible value of index where i can be moved. If
// there is no such lower bound (i.e. i doesn't depend on any previous
// instruction), this method returns zero index.
func (b *block) lowerBound(i int) int {
	ins := b.index(i)
	idx := findBound(func(i, j int) bool { return i > j },
		ins.trueDepsBack,
		ins.antiDepsBack,
		ins.outputDepsBack,
		ins.controlDepsBack,
		ins.specialDepsBack,
	)

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
	idx := findBound(func(i, j int) bool { return i < j },
		ins.trueDepsFwd,
		ins.antiDepsFwd,
		ins.outputDepsFwd,
		ins.controlDepsFwd,
		ins.specialDepsFwd,
	)

	if idx < 0 {
		return b.Len() - 1
	}
	return idx - 1
}

func (b *block) setIndex(i int) { b.idx = i }
