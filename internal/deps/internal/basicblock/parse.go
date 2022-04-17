package basicblock

import (
	"mltwist/internal/repr"
	"mltwist/pkg/model"
	"fmt"
	"sort"
)

// Parse identifies basic blocks in a sequence of program instructions.
func Parse(instrs []repr.Instruction) ([][]repr.Instruction, error) {
	sort.Slice(instrs, func(i, j int) bool {
		return instrs[i].Address < instrs[j].Address
	})

	seqs := pipeline.apply(instrs)

	blocks, err := splitByJumpTargets(seqsToBlocks(seqs))
	if err != nil {
		return nil, fmt.Errorf("cannot split blocks by jump targets: %w", err)
	}

	sequences := make([][]repr.Instruction, len(blocks))
	for i, b := range blocks {
		sequences[i] = b.seq
	}

	return sequences, nil
}

func seqsToBlocks(seqs [][]repr.Instruction) []block {
	blocks := make([]block, len(seqs))
	for i, s := range seqs {
		blocks[i] = newBlock(s)
	}

	return blocks
}

func splitByAddress(seq []repr.Instruction) [][]repr.Instruction {
	seqs := make([][]repr.Instruction, 0, 1)
	begin := 0

	for i := range seq[1:] {
		if seq[i].Address+seq[i].ByteLen == seq[i+1].Address {
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

func controlFlowInstruction(t model.Type) bool {
	return t.Jump() || t.CJump() || t.JumpDyn()
}

func splitByJumps(seq []repr.Instruction) [][]repr.Instruction {
	seqs := make([][]repr.Instruction, 0, 1)
	begin := 0

	for i, ins := range seq {
		if t := ins.Type; controlFlowInstruction(t) {
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
			for _, j := range ins.JumpTargets {
				err := blocks.split(j)
				if err != nil {
					return nil, fmt.Errorf(
						"cannot split basic block: %w", err)
				}
			}
		}
	}

	return blocks, nil
}
