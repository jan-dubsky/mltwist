package deps

import (
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"
)

func TestSpecialDeps(t *testing.T) {
	tests := []testCase{{
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
		name: "memory_order_instructions",
		ins: []*instruction{
			testInsMem([]expr.Key{"foo"}, nil),
			testInsMem(nil, []expr.Key{"foo"}),
			testInsReg(1),
			testIns(model.TypeMemOrder, nil),
			testInsMem([]expr.Key{"foo"}, nil),
			testInsMem(nil, []expr.Key{"foo"}),
			testInsReg(2),
			testIns(model.TypeMemOrder, nil),
			testInsMem([]expr.Key{"foo"}, nil),
			testInsReg(3),
		},
		deps: []dep{
			{0, 3},
			{1, 3},
			{3, 4},
			{3, 5},
			{3, 7},
			{4, 7},
			{5, 7},
			{7, 8},
		},
	}, {
		name: "multiple_special_instr",
		ins: []*instruction{
			testInsReg(1),
			testInsReg(2),
			testIns(model.TypeCPUStateChange, nil),
			testInsReg(3),
			testInsReg(4, 1),
			testInsReg(7, 3, 1),
			testIns(model.TypeSyscall, nil),
			testInsReg(8, 4, 2),
		},
		deps: []dep{
			{0, 2},
			{1, 2},
			{2, 3},
			{2, 4},
			{2, 5},
			{2, 6},
			{3, 6},
			{4, 6},
			{5, 6},
			{6, 7},
		},
	}, {
		name: "mem_order_and_special",
		ins: []*instruction{
			testIns(model.TypeMemOrder, nil),
			testIns(model.TypeCPUStateChange, nil),
			testIns(model.TypeCPUStateChange, nil),
			testIns(model.TypeMemOrder, nil),
			testIns(model.TypeMemOrder, nil),
			testIns(model.TypeSyscall, nil),
			testIns(model.TypeSyscall, nil),
			testIns(model.TypeMemOrder, nil),
			testIns(model.TypeSyscall, nil),
		},
		deps: []dep{
			{0, 1},
			{1, 2},
			{2, 3},
			{2, 4},
			{2, 5},
			{3, 4},
			{3, 5},
			{4, 5},
			{5, 6},
			{6, 7},
			{6, 8},
			{7, 8},
		},
	}}

	runDepsTest(t, tests, findSpecialDeps)
}
