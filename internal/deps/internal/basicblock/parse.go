package basicblock

import (
	"fmt"
	"mltwist/internal/parser"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"sort"
)

// Parse identifies basic blocks in a sequence of program instructions. This
// function returns sorted (by increasing address) list of basic blocks where
// every basic block is represented as a sorted list of instructions.
//
// The entrypoint argument is in a way special as it's the jump target which by
// originates from outside of the program. Naturally there are as well other
// jump targets comming from outside (for example intel SYSEXIT handler in
// glibc). Unfortunately those require non trivial platform and OS knowledge to
// identify. The entrypoint is the only well-defined jump target comming from
// outside of the program.
func Parse(entrypoint model.Addr, instrList []parser.Instruction) ([][]Instruction, error) {
	instrs := convertInstructions(instrList)
	sort.Slice(instrs, func(i, j int) bool { return instrs[i].Addr < instrs[j].Addr })

	seqs := pipeline.apply(instrs)

	bs, err := splitByJumpTargets(seqsToBlocks(seqs))
	if err != nil {
		return nil, fmt.Errorf("cannot split blocks by jump targets: %w", err)
	}

	err = bs.split(entrypoint)
	if err != nil {
		return nil, fmt.Errorf("cannot create basic block at entry point: %w", err)
	}

	sequences := make([][]Instruction, len(bs))
	for i, b := range bs {
		sequences[i] = b.seq
	}

	return sequences, nil
}

func seqsToBlocks(seqs [][]Instruction) []block {
	blocks := make([]block, len(seqs))
	for i, s := range seqs {
		blocks[i] = newBlock(s)
	}

	return blocks
}

func splitByAddress(seq []Instruction) [][]Instruction {
	seqs := make([][]Instruction, 0, 1)
	begin := 0

	for i := range seq[1:] {
		if seq[i].NextAddr() == seq[i+1].Addr {
			continue
		}

		seqs = append(seqs, seq[begin:i+1])
		begin = i + 1
	}

	if end := len(seq); begin < end {
		seqs = append(seqs, seq[begin:end])
	}

	return seqs
}

func splitByJumps(seq []Instruction) [][]Instruction {
	seqs := make([][]Instruction, 0, 1)
	begin := 0

	for i, ins := range seq {
		if len(ins.JumpTargets) > 0 {
			seqs = append(seqs, seq[begin:i+1])
			begin = i + 1
		}
	}

	if end := len(seq); begin < end {
		seqs = append(seqs, seq[begin:end])
	}

	return seqs
}

func splitByJumpTargets(bs []block) (blocks, error) {
	blocks := make(blocks, len(bs))
	for i, b := range bs {
		blocks[i] = b
	}

	for _, b := range bs {
		for _, ins := range b.seq {
			for _, jumpExpr := range ins.JumpTargets {
				c, ok := jumpExpr.(expr.Const)
				if !ok {
					continue
				}

				addr, ok := expr.ConstUint[model.Addr](c)
				if !ok {
					continue
				}

				err := blocks.split(addr)
				if err != nil {
					return nil, fmt.Errorf(
						"cannot split basic block: %w", err)
				}
			}
		}
	}

	return blocks, nil
}
