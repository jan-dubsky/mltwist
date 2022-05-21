package deps

import (
	"mltwist/pkg/model"
	"testing"
)

func TestSpecialDeps(t *testing.T) {
	tests := []testCase{
		{
			name: "single_special_instr",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(2),
				testIns(model.TypeCPUStateChange, nil),
				testInsReg(3, 2, 1),
			},
			deps: []dep{
				{0, 2},
				{1, 2},
				{2, 3},
			},
		}, {
			name: "multiple_special_instr",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(2),
				testIns(model.TypeCPUStateChange, nil),
				testInsReg(3),
				testInsReg(4, 1),
				testIns(model.TypeMemOrder, nil),
				testInsReg(7, 3, 1),
				testInsReg(8, 4, 2),
				testIns(model.TypeSyscall, nil),
			},
			deps: []dep{
				{0, 2},
				{1, 2},
				{2, 3},
				{2, 4},
				{2, 5},
				{3, 5},
				{4, 5},
				{5, 6},
				{5, 7},
				{5, 8},
				{6, 8},
				{7, 8},
			},
		},
	}

	runDepsTest(t, tests, findSpecialDeps)
}
