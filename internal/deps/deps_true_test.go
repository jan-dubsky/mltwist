package deps

import (
	"decomp/internal/repr"
	"decomp/pkg/model"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

const regInvalid model.Register = math.MaxUint64

func testIns(ins repr.Instruction) *instruction {
	return &instruction{
		trueDepsFwd:     make(insSet, numRegs),
		trueDepsBack:    make(insSet, numRegs),
		antiDepsFwd:     make(insSet, numRegs),
		antiDepsBack:    make(insSet, numRegs),
		outputDepsFwd:   make(insSet, numRegs),
		outputDepsBack:  make(insSet, numRegs),
		controlDepsFwd:  make(insSet, 1),
		controlDepsBack: make(insSet, 0),
		specialDepsFwd:  make(insSet, 0),
		specialDepsBack: make(insSet, 0),
		blockIdx:        -1,

		Instr: ins,
	}
}

func testInsReg(out model.Register, in ...model.Register) *instruction {
	return testIns(testInsReprReg(out, in...))
}

func testInsReprReg(out model.Register, in ...model.Register) repr.Instruction {
	inRegs := make(map[model.Register]struct{}, len(in))
	for _, r := range in {
		if _, ok := inRegs[r]; ok {
			panic(fmt.Sprintf("duplicit register: %d", r))
		}

		inRegs[r] = struct{}{}
	}

	outRegs := make(map[model.Register]struct{}, 1)
	if out != regInvalid {
		outRegs[out] = struct{}{}
	}

	return repr.Instruction{
		Instruction: model.Instruction{
			InputRegistry:  inRegs,
			OutputRegistry: outRegs,
		},
	}
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
	fwdF depFunc,
	backF depFunc,
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
				t.Logf("Forward deps %d: %d\n", i, len(fwdF(ins)))
			}
			t.Logf("\n")
			for i, ins := range tt.ins {
				t.Logf("Backward deps %d: %d\n", i, len(backF(ins)))
			}

			for d := range deps {
				src := tt.ins[d.src]
				dst := tt.ins[d.dst]
				r.Contains(fwdF(src), dst, "dependency: %v", d)
				r.Contains(backF(dst), src, "dependency: %v", d)
			}

			srcs := depCnts(deps, func(d dep) int { return d.src })
			for i := range tt.ins {
				r.Len(fwdF(tt.ins[i]), srcs[i], "instruction: %d", i)
			}

			dsts := depCnts(deps, func(d dep) int { return d.dst })
			for i := range tt.ins {
				r.Len(backF(tt.ins[i]), dsts[i], "instruction: %d", i)
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

	runDepsTest(t, tests, processTrueDeps,
		func(i *instruction) insSet { return i.trueDepsFwd },
		func(i *instruction) insSet { return i.trueDepsBack },
	)
}
