package deps

import (
	"mltwist/internal/deps/internal/basicblock"
	"mltwist/pkg/expr"
	"testing"
)

func testInsJumpTarget(exprs ...expr.Expr) *instruction {
	return testIns(basicblock.Instruction{JumpTargets: exprs})
}

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
				testInsJumpTarget(
					expr.NewConstUint[uint8](56, expr.Width32),
				),
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
				testInsJumpTarget(expr.NewCond(expr.Eq,
					expr.Zero,
					expr.One,
					expr.Zero,
					expr.One,
					expr.Width32,
				)),
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
				testInsJumpTarget(expr.NewRegLoad("r1", expr.Width32)),
			},
			deps: []dep{
				{0, 3},
				{1, 3},
				{2, 3},
			},
		},
	}

	runDepsTest(t, tests, processControlDeps)
}
