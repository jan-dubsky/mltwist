package deps

import (
	"mltwist/pkg/expr"
	"testing"
)

func TestOutputDeps_Register(t *testing.T) {
	tests := []testCase{{
		name: "simple_add",
		ins: []*instruction{
			testInsReg(1),
			testInsReg(2),
			testInsReg(2, 2, 1),
			testInsReg(3),
		},
		deps: []dep{
			{1, 2},
		},
	}, {
		name: "add_multiple",
		ins: []*instruction{
			testInsReg(1, 1), // Cannot be output dependent in itself.
			testInsReg(2),
			testInsReg(3, 2, 1),
			testInsReg(4, 4),
			testInsReg(5, 4, 3),
			testInsReg(3, 5),
			testInsReg(4),
			testInsReg(2, 4, 3),
		},
		deps: []dep{
			{2, 5},
			{3, 6},
			{1, 7},
		},
	}, {
		name: "no_deps",
		ins: []*instruction{
			testInsReg(1, 1),
			testInsReg(3, 3),
			testInsReg(regInvalid),
			testInsReg(2, 3, 1),
			testInsReg(4, 3, 2),
			testInsReg(5, 4, 3),
			testInsReg(regInvalid),
			testInsReg(6, 2, 3),
		},
		deps: nil,
	}}

	runDepsTest(t, tests, findOutputDeps)
}

func TestOutputDeps_Memory(t *testing.T) {
	tests := []testCase{{
		name: "read_write_single",
		ins: []*instruction{
			testInsMem(nil, nil),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, nil),
		},
		deps: []dep{
			{1, 3},
		},
	}, {
		name: "read_write_chain_single_memory",
		ins: []*instruction{
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
			testInsMem(nil, nil),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, nil),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
		},
		deps: []dep{
			{1, 3},
			{3, 6},
			{6, 7},
		},
	}, {
		name: "multiple_memories",
		ins: []*instruction{
			testInsMem([]expr.Key{"m2"}, []expr.Key{"m2"}),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m2"}),
			testInsMem([]expr.Key{"m2"}, nil),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m2", "m1"}),
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem(nil, []expr.Key{"m2"}),
			testInsMem([]expr.Key{"m1"}, nil),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m2"}, []expr.Key{"m1", "m2"}),
		},
		deps: []dep{
			// m1
			{1, 3},
			{3, 6},
			{6, 7},
			// m2
			{0, 2},
			{2, 8},
		},
	}}

	runDepsTest(t, tests, findOutputDeps)
}
