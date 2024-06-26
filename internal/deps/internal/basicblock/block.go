package basicblock

import (
	"fmt"
	"mltwist/pkg/model"
	"sort"
)

// block represents a sequence of instructions which is always executed in the
// code sequentially (i.e. there are no jump instructions and jump targets
// inside the basic block). Such instruction sequence is typically referred as
// basic block in compilers.
//
// As most of known CPU architectures allow jump instructions to a dynamic value
// (register value), it's not always possible to identify all basic blocks. This
// implies that a dynamic jump instructions might jump into the middle of the
// basic block just because it't not possible to identify jump target of this
// dynamic jump during the decompilation process.
type block[T Instruction] struct {
	seq    []T
	length model.Addr
}

func newBlock[T Instruction](seq []T) block[T] {
	return block[T]{
		seq:    seq,
		length: seqBytes(seq),
	}
}

// seqBytes calculates sum of instruction lengths in a sequence.
func seqBytes[T Instruction](seq []T) model.Addr {
	var length model.Addr
	for _, ins := range seq {
		length += ins.End() - ins.Begin()
	}
	return length
}

// begin returns inclusive start address of b.
func (b block[T]) begin() model.Addr { return b.seq[0].Begin() }

// end returns exclusive end address of b.
func (b block[T]) end() model.Addr { return b.begin() + b.length }

// Containts check if addr is inside the basic block.
func (b block[T]) contains(addr model.Addr) bool {
	return addr >= b.begin() && addr < b.end()
}

// split creates 2 new basic blocks consisting of instructions of b, but
// separated by addr respectively. The instruction starting at addr will be
// already included in the later block. If addr doesn't belong to the block or
// is not at an instruction boundary, this method returns an error.
//
// Please note that even though adds==b.Begin() is technically correct and will
// result in empty first block returned, it makes just little sense to perform
// split at b.Begin() address.
func (b block[T]) split(addr model.Addr) (block[T], block[T], error) {
	if !b.contains(addr) {
		err := fmt.Errorf("block doesn't contain address 0x%x", addr)
		return block[T]{}, block[T]{}, err
	}

	i := sort.Search(len(b.seq), func(i int) bool { return b.seq[i].Begin() >= addr })
	if i == len(b.seq) {
		err := fmt.Errorf("instruction at address 0x%x not found", addr)
		return block[T]{}, block[T]{}, err
	}
	if b.seq[i].Begin() != addr {
		err := fmt.Errorf("address 0x%x is not at instruction boundary", addr)
		return block[T]{}, block[T]{}, err
	}

	return newBlock(b.seq[:i]), newBlock(b.seq[i:]), nil
}

// blocks is an ordered sequence of basic blocks which allows fast and efficient
// splitting of basic blocks at a given address. The order of basic blocks is an
// ascending order of memory addresses. As basic blocks are non-overlapping in
// memory, the ordering applies to both Begin() and End() of the block.
type blocks[T Instruction] []block[T]

// split splits block containing addr into 2 blocks using Split(addr) and
// modifies blocks to contain both new blocks instead of the one splitted.
func (bs *blocks[T]) split(addr model.Addr) error {
	idx := sort.Search(len(*bs), func(i int) bool { return (*bs)[i].end() > addr })
	if idx == len(*bs) || addr < (*bs)[idx].begin() {
		return fmt.Errorf("no basic block with address 0x%x found", addr)
	}

	// Address is at basic block start, so no splitting is necessary.
	if (*bs)[idx].begin() == addr {
		return nil
	}

	b1, b2, err := (*bs)[idx].split(addr)
	if err != nil {
		return fmt.Errorf("cannot split block: %w", err)
	}

	*bs = append(*bs, block[T]{})
	for i := len(*bs) - 1; i >= idx+2; i-- {
		(*bs)[i] = (*bs)[i-1]
	}

	(*bs)[idx], (*bs)[idx+1] = b1, b2
	return nil
}
