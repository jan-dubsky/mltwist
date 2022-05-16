package basicblock

type splitFunc[T Instruction] func(seq []T) [][]T

func pipelineApply[T Instruction](seqs [][]T, fs ...splitFunc[T]) [][]T {
	for _, f := range fs {
		// We have no clue how many new blocks will this stage create,
		// but we know the lower bound, so we pre-allocate the lower
		// bound.
		newSeqs := make([][]T, 0, len(seqs))
		for _, b := range seqs {
			newSeqs = append(newSeqs, f(b)...)
		}

		seqs = newSeqs
	}

	return seqs
}
