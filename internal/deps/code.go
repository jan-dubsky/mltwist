package deps

import (
	"fmt"
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/parser"
	"mltwist/pkg/model"
	"sort"
)

// Code represents analyzed set of instructions from a single program which are
// split into basic blocks and have dependencies in between one another.
//
// The program code is represented as a list of basic blocks where every basic
// block contains some instructions. Dependencies in between instructions are
// then tracked only within the basic block.
type Code struct {
	entrypoint model.Addr
	blocks     []*block

	// blocksByAddr is list of blocks sorted by their ascending begin
	// addresses. This list is used for O(log(n)) search in blocks based on
	// the address. Blocks field cannot be used for such a purpose as it can
	// be modified by Move method.
	blocksByAddr []*block

	// instrCnt is number of instructions in all basic blocks.
	instrCnt int
}

// NewCode finds basic blocks in the program and identifies instruction
// dependencies within basic blocks.
func NewCode(entrypoint model.Addr, seq []parser.Instruction) (*Code, error) {
	ins := make([]*instruction, len(seq))
	for i, instruction := range seq {
		ins[i] = newInstruction(instruction)
	}

	seqs, err := basicblock.Parse(entrypoint, ins)
	if err != nil {
		return nil, fmt.Errorf("cannot find basic blocks: %w", err)
	}

	blocks := make([]*block, len(seqs))
	for i, seq := range seqs {
		blocks[i] = newBlock(i, seq)
	}

	var instrCnt int
	for _, b := range blocks {
		instrCnt += b.Num()
	}

	blocksByAddr := make([]*block, len(blocks))
	copy(blocksByAddr, blocks)

	return &Code{
		entrypoint:   entrypoint,
		blocks:       blocks,
		blocksByAddr: blocksByAddr,

		instrCnt: instrCnt,
	}, nil
}

// Entrypoint returns address of program entrypoint.
func (c *Code) Entrypoint() model.Addr { return c.entrypoint }

// Len returns number of basic blocks in the program.
func (c *Code) Len() int { return len(c.blocks) }

// Index returns ith basic block in the program.
//
// This method panics for negative values of i as well as for i greater or equal
// to p.Len().
func (c *Code) Index(i int) Block { return wrapBlock(c.blocks[i]) }

// Blocks returns list of all basic blocks in the program. Caller is allowed to
// modify the returned array.
//
// This function allocates a new array of Blocks, so its cost is O(n). In case
// you need to access just a few blocks (not all), prefer using p.Index method.
func (c *Code) Blocks() []Block {
	blocks := make([]Block, len(c.blocks))
	for i, b := range c.blocks {
		blocks[i] = wrapBlock(b)
	}
	return blocks
}

// NumInstr counts number of instructions in all all basic blocks in the code.
func (c *Code) NumInstr() int { return c.instrCnt }

// Move moves basic block at index from to index to.
//
// Unlike for instructions within a single basic block, move of a basic block
// has no effect on it's address in the program. For detailed explanation why we
// don't move basic blocks and which effect such a move could have to relative
// jump instructions, please read doc-comment of Move method of Block in this
// package.
func (c *Code) Move(from int, to int) error {
	if err := checkFromToIndex(from, to, len(c.blocks)); err != nil {
		return fmt.Errorf("cannot move %d to %d: %w", from, to, err)
	}

	move(c.blocks, from, to)
	return nil
}

// Address find a block containing address a. If no such block exists in the
// code, this method returns zero value of block and false.
func (c *Code) Address(a model.Addr) (Block, bool) {
	i := sort.Search(len(c.blocksByAddr), func(i int) bool {
		return c.blocksByAddr[i].end > a
	})
	if i == len(c.blocksByAddr) {
		return Block{}, false
	}

	b := c.blocksByAddr[i]
	if b.begin > a {
		return Block{}, false
	}

	return wrapBlock(b), true
}
