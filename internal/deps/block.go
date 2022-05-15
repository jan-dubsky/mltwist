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

	// seq is list of instruction sorted by their ascending in-memory
	// addresses.
	seq []*instruction

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

// Len returns length of the block in bytes.
func (b *block) Len() model.Addr { return b.end - b.begin }

// End returns in-memory address of the first byte behind the block. The address
// relates to the original address space of a binary.
func (b *block) End() model.Addr { return b.end }

// Num returns number of instructions in b.
func (b *block) Num() int { return len(b.seq) }

// Idx returns index of an instruction in list of basic blocks.
func (b *block) Idx() int { return b.idx }

// index returns instruction at index i.
//
// This method is non-exported as it returns internal object. Expo
func (b *block) index(i int) *instruction { return b.seq[i] }

// Move moves instruction in the block from index from to index to. All
// instructions in between from and to are shifted one instruction back or
// forward respectively. This method will fail in case the move violates any
// instruction dependency constraints or if either from or to are not valid
// indices of an instruction in the block.
//
// Move of an instruction correctly changes its address, but instruction bytes
// are left untouched. This might result in a state when bytes or relative jumps
// (jumps to an instruction address plus some offset) contain different jump
// target than the real instruction.
//
// Unfortunately it's not possible to modify instruction bytes to be sure that
// the instruction binary encoding is always valid even after the move
// operation. To do so, we'd need a platform-dependent module which would encode
// the change into instruction bytes. But even if we had such a module, the
// module might not be able to assemble a valid instruction opcode for the
// target architecture. In a typical CPU architecture a relative jump can jump
// only in some small area around the instruction (typically a few MB). For
// longer jumps, the CPU architecture typically requires to load the value into
// a register and jump to an address in register. Consequently, a move of a
// single instruction might result in expansion of a single instruction into
// multiple. Moreover, remember that we move whole range of [begin, end] or
// [end, begin] respectively. So in the worst case there can be multiple
// instruction expanded. This might result in several paradoxes.
//
// The first paradox or more a suboptimal situation would be that expansion of
// an instruction would result in basic block size growth. As basic block can
// have adjacent basic blocks both before and after, to grow a basic block, we'd
// need to shift the basic blocks following in the program address space. But
// such a move would trigger other relative jump instruction modification, which
// could again result in expansion and basic block move. So in the worst case
// (more an absurd case) scenario, one instruction move could cause expansion of
// all relative jumps in the program followed by move of all basic blocks.
//
// Another paradox is that by moving basic blocks, we'd have to change not just
// relative jumps, but also absolute jumps. Let's assume that absolute jump
// always jumps to an address in register. Such an assumption doesn't cause any
// loss of generality. The problem with jumps to registry is that we have to
// somehow will the register with the value. This in generic case requires to
// perform inter-block analysis as the register in general doesn't have to be
// filled by a constant just before the jump instruction. It's obvious why
// something as inter-block analysis of all possible values in unacceptable. But
// even the where constant value is loaded into a register just before the jump
// instruction is non-trivial. The problem is that first, we'd still need to
// perform in-block analysis of the data flow. But second, some architectures do
// not have instructions to load the whole it width of a register in a single
// instruction. And in such a case, we might again face the situation that the
// kind of constant load used in the code is not able to load so big constant
// value. So again, we'd have to expand single constant load instruction into
// multiple instructions.
//
// The third paradox is that we still don't know addresses of some jumps. It's
// true that we ignore those jumps when we move instruction inside the basic
// block and we cannot be sure that one of those won't split our basic block
// into two. But there is still quite a significant difference in between moving
// instructions and moving blocks. Even though compilers nowadays are very
// sophisticated, they still benefit from basic blocks having fixed boundaries
// (as this allows optimizations). So it's reasonable to assume that jumps won't
// very often jump into the middle of basic block. On the other hand it's very
// likely that some of those jumps we don't know there they jump will jump to
// the beginning of a basic block. So if we move a basic block, we will most
// likely change the semantics of the program remarkably. It's also worth
// mentioning that in case of instruction move, the user is well aware of what
// he/she does and it able to manually check (or make a conscious assumption)
// that such instruction move won't change program semantics due to unknown
// jumps. On the other hand if a single instruction move results in move of K
// blocks and expansion of L instructions, the users ability to validate
// correctness of such a step is very limited.
//
// For all the reasons described above, move of an instruction doesn't change
// its bytes even for the cost of making the instruction opcode incorrect in the
// modified program code.
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
		if u := b.UpperBound(from); u < to {
			return fmt.Errorf("upper bound for move is: %d", u)
		}
	} else if from > to {
		if l := b.LowerBound(from); l > to {
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

// LowerBound finds the lowest possible value of index where i can be moved. If
// there is no such lower bound (i.e. i doesn't depend on any previous
// instruction), this method returns zero index.
func (b *block) LowerBound(i int) int {
	ins := b.index(i)
	idx := findBound(func(i, j int) bool { return i > j }, ins.depsBack)

	if idx < 0 {
		return 0
	}
	return idx + 1
}

// UpperBound finds the highest possible of index where i can be moved. If there
// is no such upper bound (i.e. i doesn't depend on any later instruction), this
// method returns b.Len() - 1.
func (b *block) UpperBound(i int) int {
	ins := b.index(i)
	idx := findBound(func(i, j int) bool { return i < j }, ins.depsFwd)

	if idx < 0 {
		return b.Num() - 1
	}
	return idx - 1
}

// setAddr is an empty implementation of address setter which allows us to use
// the same algorithm for both instructions and blocks. As blocks are not
// allowed to move in memory. No other (than empty) implementation of this
// function makes sense then.
func (b *block) setAddr(_ model.Addr) {}
func (b *block) setIndex(i int)       { b.idx = i }

// Address finds an instruction with address a in the block. If a is not in the
// block or a is not start of an instruction (is in the middle of an
// instruction), this function returns zero value of instruction and false.
func (b *block) Address(a model.Addr) (Instruction, bool) {
	i := sort.Search(len(b.seq), func(i int) bool {
		return b.seq[i].Begin() >= a
	})
	if i == len(b.seq) {
		return Instruction{}, false
	}

	ins := b.seq[i]
	if ins.Begin() != a {
		return Instruction{}, false
	}

	return wrapInstruction(ins), true
}
