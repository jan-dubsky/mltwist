package deps

import (
	"fmt"
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/internal/parser"
	"mltwist/pkg/model"
	"sort"
)

type Program struct {
	entrypoint   model.Addr
	blocks       []*block
	blocksByAddr []*block
}

func NewProgram(entrypoint model.Addr, seq []parser.Instruction) (*Program, error) {
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

	return &Program{
		entrypoint:   entrypoint,
		blocks:       blocks,
		blocksByAddr: blocksByAddr,
	}, nil
}

func (p *Program) Entrypoint() model.Addr { return p.entrypoint }

// Len returns number of basic blocks in the program.
func (p *Program) Len() int { return len(p.blocks) }

// Index returns ith basic block in the program.
//
// This method panics for negative values of i as well as for i greater or equal
// to p.Len().
func (p *Program) Index(i int) Block { return wrapBlock(p.blocks[i]) }

// Blocks returns list of all basic blocks in the program. Caller is allowed to
// modify the returned array.
//
// This function allocates a new array of Blocks, so its cost is O(n). In case
// you need to access just a few blocks (not all), prefer using p.Index method.
func (p *Program) Blocks() []Block {
	blocks := make([]Block, len(p.blocks))
	for i, b := range p.blocks {
		blocks[i] = wrapBlock(b)
	}
	return blocks
}

// NumInstr counts number of instructions in all all basic blocks in the
// program.
func (p *Program) NumInstr() int {
	var instrs int
	for _, b := range p.blocks {
		instrs += b.Len()
	}

	return instrs
}

// Move moves basic block at index from to index to.
func (p *Program) Move(from int, to int) error {
	if err := checkFromToIndex(from, to, len(p.blocks)); err != nil {
		return fmt.Errorf("cannot move %d to %d: %w", from, to, err)
	}

	move(p.blocks, from, to)
	return nil
}

func (p *Program) Addr(a model.Addr) (Block, bool) {
	i := sort.Search(len(p.blocksByAddr), func(i int) bool {
		return p.blocksByAddr[i].end > a
	})
	if i == len(p.blocksByAddr) {
		return Block{}, false
	}

	b := p.blocksByAddr[i]
	if b.begin > a {
		return Block{}, false
	}

	return wrapBlock(b), true
}

func (p *Program) AddrIns(a model.Addr) (Instruction, bool) {
	block, ok := p.Addr(a)
	if !ok {
		return Instruction{}, false
	}

	ins, ok := block.Addr(a)
	if !ok {
		return Instruction{}, false
	}

	return ins, true
}
