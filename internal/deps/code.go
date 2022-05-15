package deps

import (
	"fmt"
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/parser"
	"mltwist/pkg/model"
	"sort"
)

type Code struct {
	entrypoint   model.Addr
	blocks       []*block
	blocksByAddr []*block
}

func NewCode(entrypoint model.Addr, seq []parser.Instruction) (*Code, error) {
	seqs, err := basicblock.Parse(entrypoint, seq)
	if err != nil {
		return nil, fmt.Errorf("cannot find basic blocks: %w", err)
	}

	blocks := make([]*block, len(seqs))
	for i, seq := range seqs {
		blocks[i] = newBlock(i, seq)
	}

	blocksByAddr := make([]*block, len(blocks))
	copy(blocksByAddr, blocks)

	return &Code{
		entrypoint:   entrypoint,
		blocks:       blocks,
		blocksByAddr: blocksByAddr,
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

// NumInstr counts number of instructions in all all basic blocks in the
// program.
func (c *Code) NumInstr() int {
	var instrs int
	for _, b := range c.blocks {
		instrs += b.Len()
	}

	return instrs
}

// Move moves basic block at index from to index to.
func (c *Code) Move(from int, to int) error {
	if err := checkFromToIndex(from, to, len(c.blocks)); err != nil {
		return fmt.Errorf("cannot move %d to %d: %w", from, to, err)
	}

	move(c.blocks, from, to)
	return nil
}

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

func (c *Code) AddressIns(a model.Addr) (Instruction, bool) {
	block, ok := c.Address(a)
	if !ok {
		return Instruction{}, false
	}

	ins, ok := block.Address(a)
	if !ok {
		return Instruction{}, false
	}

	return ins, true
}
