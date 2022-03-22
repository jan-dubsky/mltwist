package basicblock

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
	"fmt"
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
type block struct {
	Seq    []repr.Instruction
	length model.Address
}

func newBlock(seq []repr.Instruction) block {
	return block{Seq: seq, length: seqBytes(seq)}
}

// seqBytes calculates sum of instruction lengths in a sequence.
func seqBytes(seq []repr.Instruction) model.Address {
	var length model.Address
	for _, ins := range seq {
		length += ins.ByteLen
	}
	return length
}

// begin returns inclusive start address of b.
func (b block) begin() model.Address { return b.Seq[0].Address }

// end returns exclusive end address of b.
func (b block) end() model.Address { return b.begin() + b.length }

// Containts check if addr is inside the basic block.
func (b block) contains(addr model.Address) bool {
	return addr >= b.begin() && addr < b.end()
}

// Split creates 2 new basic blocks consisting of instructions of b, but
// separated by addr respectively. The instruction starting at addr will be
// already included in the later block. If addr doesn't belong to the block or
// is not at an instruction boundary, this method returns an error.
//
// Please note that even though adds==b.Begin() is technically correct and will
// result in empty first block returned, it makes just little sense to perform
// split at b.Begin() address.
func (b block) Split(addr model.Address) (block, block, error) {
	if !b.contains(addr) {
		err := fmt.Errorf("block doesn't contain address 0x%x", addr)
		return block{}, block{}, err
	}

	i := sort.Search(len(b.Seq), func(i int) bool { return b.Seq[i].Address >= addr })
	if i == len(b.Seq) || b.Seq[i].Address != addr {
		err := fmt.Errorf("address 0x%x is not at instruction boundary", addr)
		return block{}, block{}, err
	}

	return newBlock(b.Seq[:i]), newBlock(b.Seq[i:]), nil
}

// blocks is an ordered sequence of basic blocks which allows fast and efficient
// splitting of basic blocks at a given address. The order of basic blocks is an
// ascending order of memory addresses. As basic blocks are non-overlapping in
// memory, the ordering applies to both Begin() and End() of the block.
type blocks []block

// split splits block containing addr into 2 blocks using Split(addr) and
// modifies blocks to contain both new blocks instead of the one splitted.
func (bs *blocks) split(addr model.Address) error {
	idx := sort.Search(len(*bs), func(i int) bool { return (*bs)[i].end() > addr })
	if idx == len(*bs) {
		return fmt.Errorf("no basic block with address 0x%x found", addr)
	}

	// Address is at basic block start, so no splitting is necessary.
	if (*bs)[idx].begin() == addr {
		return nil
	}

	b1, b2, err := (*bs)[idx].Split(addr)
	if err != nil {
		return fmt.Errorf("cannot split block: %w", err)
	}

	*bs = append(*bs, block{})
	for i := idx + 1; i < len(*bs)-1; i++ {
		(*bs)[i+1] = (*bs)[i]
	}

	(*bs)[idx], (*bs)[idx+1] = b1, b2
	return nil
}
