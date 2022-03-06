package deps

import (
	"testing"
)

func TestAntiDeps_Register(t *testing.T) {
	tests := []testCase{
		{
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
		},
		{
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
		},
		{
			name: "no_deps",
			ins: []*instruction{
				testInsReg(1, 2),
				testInsReg(3, 4),
				testInsReg(5, 1, 3),
				testInsReg(6, 5, 3, 1),
				testInsReg(7, 5, 6),
			},
			deps: nil,
		},
	}

	runDepsTest(t, tests, processAntiDeps,
		func(i *instruction) insSet { return i.antiDepsFwd },
		func(i *instruction) insSet { return i.antiDepsBack },
	)
}
