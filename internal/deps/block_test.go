package deps

import (
	"decomp/internal/repr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockBounds(t *testing.T) {
	type bounds struct {
		lower int
		upper int
	}
	tests := []struct {
		name   string
		ins    []repr.Instruction
		bounds map[int]bounds
	}{
		{
			name: "simple_add",
			ins: []repr.Instruction{
				testInsReprReg(1),
				testInsReprReg(2),
				testInsReprReg(3, 1, 2),
			},
			bounds: map[int]bounds{
				0: {lower: 0, upper: 1},
				1: {lower: 0, upper: 1},
				2: {lower: 2, upper: 2},
			},
		},
		{
			name: "multiple_adds",
			ins: []repr.Instruction{
				testInsReprReg(1, 1),
				testInsReprReg(3),
				testInsReprReg(2, 2),
				testInsReprReg(4, 1, 3),
				testInsReprReg(5, 2, 1),
				testInsReprReg(6, 3, 4),
				testInsReprReg(7, 1, 3),
				testInsReprReg(3),
				testInsReprReg(1),
				testInsReprReg(8),
			},
			bounds: map[int]bounds{
				0: {lower: 0, upper: 2},
				1: {lower: 0, upper: 2},
				2: {lower: 0, upper: 3},
				3: {lower: 2, upper: 4},
				4: {lower: 3, upper: 7},
				5: {lower: 4, upper: 6},
				6: {lower: 2, upper: 6},
				7: {lower: 7, upper: 9},
				8: {lower: 7, upper: 9},
				9: {lower: 0, upper: 9},
			},
		},
		{
			name: "anti_dependencies",
			ins: []repr.Instruction{
				testInsReprReg(1),
				testInsReprReg(2, 7, 2),
				testInsReprReg(3, 5),
				testInsReprReg(5, 4),
				testInsReprReg(4, 8),
				testInsReprReg(6, 9),
				testInsReprReg(7, 7),
			},
			bounds: map[int]bounds{
				0: {lower: 0, upper: 6},
				1: {lower: 0, upper: 5},
				2: {lower: 0, upper: 2},
				3: {lower: 3, upper: 3},
				4: {lower: 4, upper: 6},
				5: {lower: 0, upper: 6},
				6: {lower: 2, upper: 6},
			},
		},
	}

	idxs := func(instrs insSet) []int {
		indexes := make([]int, 0, len(instrs))
		for ins := range instrs {
			indexes = append(indexes, ins.blockIdx)
		}
		return indexes
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			block := newBlock(tt.ins)
			for i, ins := range block.seq {
				r.Equal(i, ins.blockIdx)
			}

			for i, b := range tt.bounds {
				ins := block.Idx(i)
				l, u := block.LowerBound(ins), block.UpperBound(ins)
				r.LessOrEqual(l, u)
				r.GreaterOrEqual(l, 0)
				r.Less(u, block.Len())

				t.Logf("Index: %d\n", i)
				t.Logf("\tFwd true: %v\n", idxs(ins.i.trueDepsFwd))
				t.Logf("\tFwd anti: %v\n", idxs(ins.i.antiDepsFwd))
				t.Logf("\tFwd out: %v\n", idxs(ins.i.outputDepsFwd))
				t.Logf("\tBack true: %v\n", idxs(ins.i.trueDepsBack))
				t.Logf("\tBack anti: %v\n", idxs(ins.i.antiDepsBack))
				t.Logf("\tBack out: %v\n", idxs(ins.i.outputDepsBack))

				r.Equal(b.lower, l, "Lower bound doesn't match")
				r.Equal(b.upper, u, "Upper bound doesn't match")
			}
		})
	}
}
