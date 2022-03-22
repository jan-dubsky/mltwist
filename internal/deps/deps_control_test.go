package deps

import "testing"

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
