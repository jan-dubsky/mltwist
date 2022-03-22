package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlock_Bounds(t *testing.T) {
	// Keep in mind that last instruction in a basic block has always
	// control dependency on all other instructions.

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
				7: {lower: 7, upper: 8},
				8: {lower: 7, upper: 8},
				9: {lower: 9, upper: 9},
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
				0: {lower: 0, upper: 5},
				1: {lower: 0, upper: 5},
				2: {lower: 0, upper: 2},
				3: {lower: 3, upper: 3},
				4: {lower: 4, upper: 5},
				5: {lower: 0, upper: 5},
				6: {lower: 6, upper: 6},
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

			block := newBlock(0, tt.ins)
			for i, ins := range block.seq {
				r.Equal(i, ins.blockIdx)
			}

			for i, b := range tt.bounds {
				ins := block.Index(i)
				l, u := block.LowerBound(ins), block.UpperBound(ins)
				r.LessOrEqual(l, u)
				r.GreaterOrEqual(l, 0)
				r.Less(u, block.Len())

				t.Logf("Index: %d\n", i)
				t.Logf("\tFwd true: %v\n", idxs(ins.i.trueDepsFwd))
				t.Logf("\tFwd anti: %v\n", idxs(ins.i.antiDepsFwd))
				t.Logf("\tFwd out: %v\n", idxs(ins.i.outputDepsFwd))
				t.Logf("\tFwd control: %v\n", idxs(ins.i.controlDepsFwd))
				t.Logf("\tFwd special: %v\n", idxs(ins.i.specialDepsFwd))
				t.Logf("\tBack true: %v\n", idxs(ins.i.trueDepsBack))
				t.Logf("\tBack anti: %v\n", idxs(ins.i.antiDepsBack))
				t.Logf("\tBack out: %v\n", idxs(ins.i.outputDepsBack))
				t.Logf("\tBack control: %v\n", idxs(ins.i.controlDepsBack))
				t.Logf("\tBack special: %v\n", idxs(ins.i.specialDepsBack))

				r.Equal(b.lower, l, "Lower bound doesn't match")
				r.Equal(b.upper, u, "Upper bound doesn't match")
			}
		})
	}
}

func TestBlock_Move(t *testing.T) {
	type move struct {
		from   int
		to     int
		hasErr bool
		order  []uint64
	}
	tests := []struct {
		name   string
		numIns int
		deps   []dep
		moves  []move
	}{
		{
			name:   "single_add",
			numIns: 3,
			deps: []dep{
				{0, 2},
				{1, 2},
			},
			moves: []move{
				{
					from:  0,
					to:    1,
					order: []uint64{1, 0, 2},
				}, {
					from:   1,
					to:     2,
					hasErr: true,
					order:  []uint64{1, 0, 2},
				}, {
					from:  0,
					to:    1,
					order: []uint64{0, 1, 2},
				}, {
					from:   1,
					to:     2,
					hasErr: true,
					order:  []uint64{0, 1, 2},
				}, {
					from:   0,
					to:     2,
					hasErr: true,
					order:  []uint64{0, 1, 2},
				}, {
					from:  1,
					to:    0,
					order: []uint64{1, 0, 2},
				}, {
					from:  1,
					to:    0,
					order: []uint64{0, 1, 2},
				},
			},
		},
		{
			name:   "no_deps",
			numIns: 5,
			moves: []move{
				{
					from:  0,
					to:    0,
					order: []uint64{0, 1, 2, 3, 4},
				}, {
					from:   0,
					to:     5,
					hasErr: true,
					order:  []uint64{0, 1, 2, 3, 4},
				}, {
					from:   0,
					to:     -1,
					hasErr: true,
					order:  []uint64{0, 1, 2, 3, 4},
				}, {
					from:  0,
					to:    3,
					order: []uint64{1, 2, 3, 0, 4},
				}, {
					from:  2,
					to:    4,
					order: []uint64{1, 2, 0, 4, 3},
				}, {
					from:  3,
					to:    1,
					order: []uint64{1, 4, 2, 0, 3},
				}, {
					from:  4,
					to:    0,
					order: []uint64{3, 1, 4, 2, 0},
				},
			},
		},
		{
			name:   "multiple_adds_overlapping",
			numIns: 8,
			deps: []dep{
				{0, 4},
				{1, 6},
				{2, 4},
				{5, 6},
				{4, 7},
				{6, 7},
			},
			moves: []move{
				{
					from:  3,
					to:    7,
					order: []uint64{0, 1, 2, 4, 5, 6, 7, 3},
				}, {
					from:   0,
					to:     3,
					hasErr: true, // Dependency (0) -> (4).
					order:  []uint64{0, 1, 2, 4, 5, 6, 7, 3},
				}, {
					from:   3,
					to:     0,
					hasErr: true, // Dependency (0) -> (4).
					order:  []uint64{0, 1, 2, 4, 5, 6, 7, 3},
				}, {
					from:  0,
					to:    2,
					order: []uint64{1, 2, 0, 4, 5, 6, 7, 3},
				}, {
					from:  4,
					to:    0,
					order: []uint64{5, 1, 2, 0, 4, 6, 7, 3},
				}, {
					from:  7,
					to:    2,
					order: []uint64{5, 1, 3, 2, 0, 4, 6, 7},
				}, {
					from:  1,
					to:    5,
					order: []uint64{5, 3, 2, 0, 4, 1, 6, 7},
				}, {
					from:  1,
					to:    5,
					order: []uint64{5, 2, 0, 4, 1, 3, 6, 7},
				}, {
					from:  2,
					to:    0,
					order: []uint64{0, 5, 2, 4, 1, 3, 6, 7},
				}, {
					from:  2,
					to:    1,
					order: []uint64{0, 2, 5, 4, 1, 3, 6, 7},
				}, {
					from:  3,
					to:    2,
					order: []uint64{0, 2, 4, 5, 1, 3, 6, 7},
				}, {
					from:  5,
					to:    7,
					order: []uint64{0, 2, 4, 5, 1, 6, 7, 3},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			seq := make([]*instruction, tt.numIns)
			for i := range seq {
				seq[i] = &instruction{
					blockIdx:     i,
					DynAddress:   model.Address(i),
					trueDepsFwd:  make(insSet),
					trueDepsBack: make(insSet),
				}
			}

			for _, d := range tt.deps {
				require.NoError(t, d.Validate())

				src, dst := seq[d.src], seq[d.dst]
				src.trueDepsFwd[dst] = struct{}{}
				dst.trueDepsBack[src] = struct{}{}
			}

			block := &Block{seq: seq}
			for i, m := range tt.moves {
				m := m
				t.Run(fmt.Sprintf("move_%d", i), func(t *testing.T) {
					r := require.New(t)

					err := block.Move(m.from, m.to)
					if m.hasErr {
						r.Error(err)
					} else {
						r.NoError(err)
					}

					for i, ins := range block.seq {
						r.Equal(i, ins.blockIdx)

						addr := model.Address(m.order[i])
						t.Logf("instr addr: %d\n", ins.DynAddress)
						r.Equal(addr, ins.DynAddress)
					}
				})
			}
		})
	}
}
