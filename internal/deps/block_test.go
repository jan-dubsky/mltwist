package deps

import (
	"fmt"
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlock_New(t *testing.T) {
	testInputInsLen := func(addr model.Addr, bytes model.Addr) basicblock.Instruction {
		return basicblock.Instruction{
			Addr:  addr,
			Bytes: make([]byte, bytes),
		}
	}

	tests := []struct {
		name  string
		seq   []basicblock.Instruction
		begin model.Addr
		bytes model.Addr
	}{
		{
			name: "single_add",
			seq: []basicblock.Instruction{
				testInputInsLen(56, 2),
				testInputInsLen(58, 3),
				testInputInsLen(61, 4),
			},
			begin: 56,
			bytes: 9,
		},
		{
			name: "multiple_ins",
			seq: []basicblock.Instruction{
				testInputInsLen(128, 4),
				testInputInsLen(132, 4),
				testInputInsLen(136, 4),

				testInputInsLen(140, 2),
				testInputInsLen(142, 2),
				testInputInsLen(144, 8),

				testInputInsLen(152, 2),
				testInputInsLen(154, 4),
			},
			begin: 128,
			bytes: 30,
		},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			b := newBlock(i, tt.seq)

			r.Equal(i, b.Idx())
			r.Equal(tt.begin, b.Begin())
			r.Equal(tt.bytes, b.Bytes())
			r.Equal(tt.begin+tt.bytes, b.End())
			r.Equal(len(tt.seq), b.Len())

			for i, ins := range tt.seq {
				r.Equal(ins.Addr, b.index(i).OrigAddr())
				r.Equal(ins.Bytes, b.index(i).bytes)
				r.Equal(ins.Details, b.index(i).details)
				r.Equal(ins.Effs, b.index(i).effects)
				r.Equal(ins.JumpTargets, b.index(i).jumpTargets)
				r.Equal(ins.Type, b.index(i).typ)
			}
		})
	}
}

// testInputInsCtr is counter of testInputInsReg calls which allows to generate
// unique register names to avoid random clashes in dependency analysis.
var testInputInsCtr int64

func testInputInsReg(out uint64, in ...uint64) basicblock.Instruction {
	effects := make([]expr.Effect, 0, len(in)+1)
	id := atomic.AddInt64(&testInputInsCtr, 1)

	for i, r := range in {
		key := expr.Key(strconv.FormatUint(r, 10))
		ef := expr.NewRegStore(
			expr.NewRegLoad(key, expr.Width8),
			expr.Key(fmt.Sprintf("test_register_%d_%d", id, i)),
			expr.Width8,
		)
		effects = append(effects, ef)
	}

	if out != regInvalid {
		key := expr.Key(strconv.FormatUint(out, 10))
		effects = append(effects, expr.NewRegStore(expr.Zero, key, expr.Width8))
	}

	return basicblock.Instruction{Effs: effects}
}

func TestBlock_Bounds(t *testing.T) {
	// Keep in mind that last instruction in a basic block has always
	// control dependency on all other instructions.

	type bounds struct {
		lower int
		upper int
	}
	tests := []struct {
		name   string
		ins    []basicblock.Instruction
		bounds map[int]bounds
	}{
		{
			name: "simple_add",
			ins: []basicblock.Instruction{
				testInputInsReg(1),
				testInputInsReg(2),
				testInputInsReg(3, 1, 2),
			},
			bounds: map[int]bounds{
				0: {lower: 0, upper: 1},
				1: {lower: 0, upper: 1},
				2: {lower: 2, upper: 2},
			},
		},
		{
			name: "multiple_adds",
			ins: []basicblock.Instruction{
				testInputInsReg(1, 1),
				testInputInsReg(3),
				testInputInsReg(2, 2),
				testInputInsReg(4, 1, 3),
				testInputInsReg(5, 2, 1),
				testInputInsReg(6, 3, 4),
				testInputInsReg(7, 1, 3),
				testInputInsReg(3),
				testInputInsReg(1),
				testInputInsReg(8),
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
			ins: []basicblock.Instruction{
				testInputInsReg(1),
				testInputInsReg(2, 7, 2),
				testInputInsReg(3, 5),
				testInputInsReg(5, 4),
				testInputInsReg(4, 8),
				testInputInsReg(6, 9),
				testInputInsReg(7, 7),
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

			block := newBlock(0, tt.ins)
			for i, ins := range block.seq {
				r.Equal(i, ins.blockIdx)
			}

			for i, b := range tt.bounds {
				l, u := block.lowerBound(i), block.upperBound(i)
				r.LessOrEqual(l, u)
				r.GreaterOrEqual(l, 0)
				r.Less(u, block.Len())

				ins := block.index(i)
				t.Logf("Index: %d\n", i)
				t.Logf("\tFwd: %v\n", idxs(ins.depsFwd))
				t.Logf("\tBack: %v\n", idxs(ins.depsBack))

				r.Equal(b.lower, l, "Lower bound doesn't match for %d", i)
				r.Equal(b.upper, u, "Upper bound doesn't match for %d", i)
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
					blockIdx: i,
					origAddr: model.Addr(i),
					currAddr: model.Addr(i),
					bytes:    make([]byte, 1),
					depsFwd:  make(insSet),
					depsBack: make(insSet),
				}
			}

			for _, d := range tt.deps {
				require.NoError(t, d.Validate())

				src, dst := seq[d.src], seq[d.dst]
				src.depsFwd[dst] = struct{}{}
				dst.depsBack[src] = struct{}{}
			}

			block := &block{seq: seq}
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

						addr := model.Addr(m.order[i])
						r.Equal(addr, ins.OrigAddr())
					}
				})
			}
		})
	}
}

func TestBlock_Address(t *testing.T) {
	addrIns := func(a model.Addr, id model.Addr) basicblock.Instruction {
		return basicblock.Instruction{
			Addr:  a,
			Bytes: make([]byte, id),
		}
	}

	tests := []struct {
		name  string
		block *block
		addr  model.Addr
		id    model.Addr
	}{{
		name: "first_ins",
		block: newBlock(0, []basicblock.Instruction{
			addrIns(0, 2),
			addrIns(2, 4),
			addrIns(6, 8),
			addrIns(14, 16),
		}),
		addr: 0,
		id:   2,
	}, {
		name: "last_ins",
		block: newBlock(0, []basicblock.Instruction{
			addrIns(0, 2),
			addrIns(2, 4),
			addrIns(6, 8),
			addrIns(14, 16),
		}),
		addr: 14,
		id:   16,
	}, {
		name: "middle_ins",
		block: newBlock(0, []basicblock.Instruction{
			addrIns(0, 2),
			addrIns(2, 4),
			addrIns(6, 8),
			addrIns(14, 16),
			addrIns(30, 32),
		}),
		addr: 6,
		id:   8,
	}, {
		name: "in_between_instructions",
		block: newBlock(0, []basicblock.Instruction{
			addrIns(0, 2),
			addrIns(2, 4),
			addrIns(6, 8),
			addrIns(14, 16),
			addrIns(30, 32),
		}),
		addr: 8,
		id:   model.MaxAddress,
	}, {
		name: "behind_last_addr",
		block: newBlock(0, []basicblock.Instruction{
			addrIns(0, 2),
			addrIns(2, 4),
			addrIns(6, 8),
			addrIns(14, 16),
			addrIns(30, 32),
		}),
		addr: 64,
		id:   model.MaxAddress,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			ins, ok := tt.block.Address(tt.addr)
			if tt.id == model.MaxAddress {
				r.False(ok)
				r.Zero(ins)
				return
			}

			r.True(ok)
			r.Equal(tt.id, ins.Len())
		})
	}
}
