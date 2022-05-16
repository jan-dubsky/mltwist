package basicblock

import (
	"fmt"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"sort"
)

type Instruction interface {
	Begin() model.Addr
	End() model.Addr
	Jumps() []expr.Expr
}

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
func Parse[T Instruction](entrypoint model.Addr, seq []T) ([][]T, error) {
	sort.Slice(seq, func(i, j int) bool { return seq[i].Begin() < seq[j].Begin() })

	seqs := pipelineApply([][]T{seq}, splitByAddress[T], splitByJumps[T])

	bs, err := splitByJumpTargets(seqsToBlocks(seqs))
	if err != nil {
		return nil, fmt.Errorf("cannot split blocks by jump targets: %w", err)
	}

	err = bs.split(entrypoint)
	if err != nil {
		return nil, fmt.Errorf("cannot create basic block at entry point: %w", err)
	}

	sequences := make([][]T, len(bs))
	for i, b := range bs {
		sequences[i] = b.seq
	}

	return sequences, nil
}

func seqsToBlocks[T Instruction](seqs [][]T) []block[T] {
	blocks := make([]block[T], len(seqs))
	for i, s := range seqs {
		blocks[i] = newBlock(s)
	}

	return blocks
}

func splitByAddress[T Instruction](seq []T) [][]T {
	seqs := make([][]T, 0, 1)
	begin := 0

	for i := range seq[1:] {
		if seq[i].End() == seq[i+1].Begin() {
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

func splitByJumps[T Instruction](seq []T) [][]T {
	seqs := make([][]T, 0, 1)
	begin := 0

	for i, ins := range seq {
		if len(ins.Jumps()) > 0 {
			seqs = append(seqs, seq[begin:i+1])
			begin = i + 1
		}
	}

	if end := len(seq); begin < end {
		seqs = append(seqs, seq[begin:end])
	}

	return seqs
}

func splitByJumpTargets[T Instruction](bs blocks[T]) (blocks[T], error) {
	blocks := make(blocks[T], len(bs))
	copy(blocks, bs)

	for _, b := range bs {
		for _, ins := range b.seq {
			for _, jumpExpr := range ins.Jumps() {
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
