package deps

import (
	"decomp/pkg/model"
	"testing"
)

func TestControlDeps(t *testing.T) {
	tests := []testCase{
		{
			name: "simple_add",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(2),
				testInsReg(3, 2, 1),
				testInsReg(0),
			},
		},
		{
			name: "last_is_jump",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(2),
				testInsReg(3, 2, 1),
				testTypeIns(model.TypeJump),
			},
			deps: []dep{
				{0, 3},
				{1, 3},
				{2, 3},
			},
		},
		{
			name: "last_is_conditional_jump",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(2),
				testInsReg(3, 2, 1),
				testTypeIns(model.TypeCJump),
			},
			deps: []dep{
				{0, 3},
				{1, 3},
				{2, 3},
			},
		},
		{
			name: "last_is_dyn_jump",
			ins: []*instruction{
				testInsReg(1),
				testInsReg(2),
				testInsReg(3, 2, 1),
				testTypeIns(model.TypeJumpDyn),
			},
			deps: []dep{
				{0, 3},
				{1, 3},
				{2, 3},
			},
		},
	}

	runDepsTest(t, tests, processControlDeps,
		func(i *instruction) insSet { return i.controlDepsFwd },
		func(i *instruction) insSet { return i.controlDepsBack },
	)
}
