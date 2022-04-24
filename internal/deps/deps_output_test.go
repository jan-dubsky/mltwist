package deps

import (
	"testing"
)

func TestOutputDeps_Register(t *testing.T) {
	tests := []testCase{
		{
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
		},
		{
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
		},
		{
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
		},
	}

	runDepsTest(t, tests, processOutputDeps)
}
