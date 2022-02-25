package basicblock

import "decomp/internal/instruction"

type splitFunc func(seq []instruction.Instruction) [][]instruction.Instruction

type splitPipeline struct {
	fs []splitFunc
}

func newPipeline(fs ...splitFunc) *splitPipeline {
	return &splitPipeline{fs: fs}
}

func (p *splitPipeline) apply(
	seqs ...[]instruction.Instruction,
) [][]instruction.Instruction {
	for _, f := range p.fs {
		// We have no clue how many new blocks will this stage create,
		// but we know the lower bound, so we pre-allocate the lower
		// bound.
		newSeqs := make([][]instruction.Instruction, 0, len(seqs))
		for _, b := range seqs {
			newSeqs = append(newSeqs, f(b)...)
		}

		seqs = newSeqs
	}

	return seqs
}

var pipeline = newPipeline(
	splitByAddress,
	splitByJumps,
)
