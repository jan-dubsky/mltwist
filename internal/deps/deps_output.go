package deps

import "mltwist/pkg/expr"

func processOutputDeps(instrs []*instruction) {
	regs := make(map[expr.Key]*instruction, numRegs)
	memory := make(map[expr.Key]*instruction, 1)

	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		processOutputDepsReg(ins, regs)
		processOutputDepsMemory(ins, memory)
	}
}

func processOutputDepsReg(ins *instruction, regs keyInsMap) {
	for r := range ins.outRegs {
		dep, ok := regs[r]
		if !ok {
			regs[r] = ins
			continue
		}

		// We are certain that i != ins.
		addDep(ins, dep)
	}
}

func processOutputDepsMemory(ins *instruction, memory keyInsMap) {
	for _, l := range ins.stores {
		dep, ok := memory[l.Key()]
		if !ok {
			continue
		}

		addDep(ins, dep)
	}

	for _, s := range ins.stores {
		memory[s.Key()] = ins
	}
}
