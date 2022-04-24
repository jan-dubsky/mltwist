package basicblock

import (
	"fmt"
	"mltwist/internal/parser"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"sort"
)

// Parse identifies basic blocks in a sequence of program instructions.
func Parse(instrList []parser.Instruction) ([][]Instruction, error) {
	instrs := convertInstructions(instrList)
	sort.Slice(instrs, func(i, j int) bool { return instrs[i].Addr < instrs[j].Addr })

	seqs := pipeline.apply(instrs)

	blocks, err := splitByJumpTargets(seqsToBlocks(seqs))
	if err != nil {
		return nil, fmt.Errorf("cannot split blocks by jump targets: %w", err)
	}

	sequences := make([][]Instruction, len(blocks))
	for i, b := range blocks {
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

func isJumpInstr(ins Instruction) bool {
	for _, e := range ins.JumpTargets {
		// Exclude those jumps which provably always jump to a following
		// instruction - There doesn't seem to be any reason for such
		// jumps, but they exist in real codes (for example in grep
		// compiled for riscv64).
		if c, ok := e.(expr.Const); ok {
			addr, ok := expr.ConstUint[model.Addr](c)
			if ok && addr == ins.NextAddr() {
				continue
			}
		}

		return true
	}

	return false
}

func splitByJumps(seq []Instruction) [][]Instruction {
	seqs := make([][]Instruction, 0, 1)
	begin := 0

	for i, ins := range seq {
		if isJumpInstr(ins) {
			seqs = append(seqs, seq[begin:i+1])
			begin = i + 1
		}
	}

	if end := len(seq); begin < end {
		seqs = append(seqs, seq[begin:end])
	}

	return seqs
}

func splitByJumpTargets(bs []block) ([]block, error) {
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
