package deps

import (
	"mltwist/pkg/expr"
	"testing"
)

func TestAntiDeps_Register(t *testing.T) {
	tests := []testCase{{
		name: "simple_add",
		ins: []*instruction{
			testInsReg(1),
			testInsReg(2),
			testInsReg(3, 2, 1),
			testInsReg(2),
		},
		deps: []dep{
			{2, 3},
		},
	}, {
		name: "add_sequential",
		ins: []*instruction{
			testInsReg(1, 1), // Cannot be antidependent on itself.
			testInsReg(2),
			testInsReg(3, 1, 2),
			testInsReg(1, 2),
			testInsReg(2, 3, 1),
			testInsReg(2, 2, 1),
		},
		deps: []dep{
			{2, 3},
			{2, 4},
			{3, 4},
		},
	}, {
		name: "no_deps",
		ins: []*instruction{
			testInsReg(1, 2),
			testInsReg(3, 4),
			testInsReg(5, 1, 3),
			testInsReg(6, 5, 3, 1),
			testInsReg(7, 5, 6),
		},
		deps: nil,
	}}

	runDepsTest(t, tests, processAntiDeps)
}

func TestAntiDeps_Memory(t *testing.T) {
	tests := []testCase{{
		name: "read_write_single",
		ins: []*instruction{
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, nil),
		},
		deps: []dep{
			{0, 1},
			{2, 3},
		},
	}, {
		name: "read_write_chain_single_memory",
		ins: []*instruction{
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, nil),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem(nil, []expr.Key{"m1"}),
			testInsMem([]expr.Key{"m1"}, nil),
			testInsMem([]expr.Key{"m1"}, []expr.Key{"m1"}),
		},
		deps: []dep{
			{0, 1},
			{4, 6},
			{5, 6},
		},
	}, {
		name: "multiple_memories",
		ins: []*instruction{
			testInsMem(nil, []expr.Key{"m2"}),
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
			{0, 2},
			{1, 2},
			{3, 8},
			{4, 6},
			{5, 8},
		},
	}}

	runDepsTest(t, tests, processAntiDeps)
}
