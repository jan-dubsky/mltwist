package deps

import (
	"fmt"
	"math"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const regInvalid uint64 = math.MaxUint64

func testIns(t model.Type, jumps []expr.Expr) *instruction {
	return &instruction{
		depsFwd:     make(insSet, numRegs),
		depsBack:    make(insSet, numRegs),
		blockIdx:    -1,
		typ:         t,
		jumpTargets: jumps,
	}
}

func testInsReg(out uint64, in ...uint64) *instruction {
	inRegs := make(regSet, len(in))
	for _, r := range in {
		rStr := strconv.FormatUint(r, 10)
		if _, ok := inRegs[rStr]; ok {
			panic(fmt.Sprintf("duplicit register: %d", r))
		}

		inRegs[rStr] = struct{}{}
	}

	ins := testIns(model.TypeNone, nil)
	ins.inRegs = inRegs
	if out != regInvalid {
		ins.outRegs = regSet{strconv.FormatUint(out, 10): struct{}{}}
	}

	return ins
}

type dep struct {
	src int
	dst int
}

func (d dep) Validate() error {
	f, t := d.src, d.dst
	if f == t {
		return fmt.Errorf("instruction cannot depend on itself: %d -> %d", f, t)
	}
	if f > t {
		return fmt.Errorf("invalid direction of dependency: %d -> %d", f, t)
	}

	return nil
}

func depMap(deps []dep) map[dep]struct{} {
	m := make(map[dep]struct{}, len(deps))
	for _, d := range deps {
		if _, ok := m[d]; ok {
			panic(fmt.Sprintf("duplicate dependency description: %v", d))
		}

		m[d] = struct{}{}
	}

	return m
}

func depCnts(deps map[dep]struct{}, f func(d dep) int) map[int]int {
	cnts := make(map[int]int)
	for d := range deps {
		cnts[f(d)]++
	}

	return cnts
}

type testCase struct {
	name string
	ins  []*instruction
	deps []dep
}

type depFunc func(i *instruction) insSet

func runDepsTest(
	t *testing.T,
	tests []testCase,
	f func(instrs []*instruction),
) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			for _, d := range tt.deps {
				r.NoError(d.Validate())
			}
			deps := depMap(tt.deps)

			f(tt.ins)

			for i, ins := range tt.ins {
				t.Logf("Forward deps %d: %d\n", i, len(ins.depsFwd))
			}
			t.Logf("\n")
			for i, ins := range tt.ins {
				t.Logf("Backward deps %d: %d\n", i, len(ins.depsBack))
			}

			for d := range deps {
				src := tt.ins[d.src]
				dst := tt.ins[d.dst]
				r.Contains(src.depsFwd, dst, "dependency: %v", d)
				r.Contains(dst.depsBack, src, "dependency: %v", d)
			}

			srcs := depCnts(deps, func(d dep) int { return d.src })
			for i := range tt.ins {
				r.Len(tt.ins[i].depsFwd, srcs[i], "instruction: %d", i)
			}

			dsts := depCnts(deps, func(d dep) int { return d.dst })
			for i := range tt.ins {
				r.Len(tt.ins[i].depsBack, dsts[i], "instruction: %d", i)
			}
		})
	}
}

func TestTrueDeps_Register(t *testing.T) {
	tests := []testCase{
		{
			name: "simple_add",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(2),
				testInsReg(3, 1, 2),
			},
			deps: []dep{
				{0, 2},
				{1, 2},
			},
		},
		{
			name: "add_overlayed",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(11),
				testInsReg(2, 3, 4),
				testInsReg(12, 13, 14),
				testInsReg(3, 1, 2),
				testInsReg(13, 11, 12),
				testInsReg(4, 3, 13),
			},
			deps: []dep{
				{0, 4},
				{2, 4},

				{1, 5},
				{3, 5},

				{4, 6},
				{5, 6},
			},
		},
		{
			name: "chain_deps",
			ins: []*instruction{
				testInsReg(2, 1),
				testInsReg(3, 2),
				testInsReg(4, 3),
				testInsReg(5, 4, 3),
				testInsReg(6, 5),
				testInsReg(7, 6, 5, 2, 1),
			},
			deps: []dep{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
				{4, 5},

				{0, 5},
				{1, 3},
				{3, 5},
			},
		},
		{
			name: "no_deps",
			ins: []*instruction{
				testInsReg(2, 3, 5),
				testInsReg(11, 15, 17, 18),
				testInsReg(1, 3, 5, 6),
				testInsReg(regInvalid, 4, 5, 6),
				testInsReg(8, 9, 0),
			},
			deps: nil,
		},
	}

	runDepsTest(t, tests, processTrueDeps)
}
