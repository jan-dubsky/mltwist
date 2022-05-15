package deps

import "mltwist/pkg/expr"

type keyInsMap map[expr.Key]*instruction

func processAntiDeps(instrs []*instruction) {
	regs := make(keyInsMap, numRegs)
	memory := make(keyInsMap, 1)

	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		processAntiDepsReg(ins, regs)
		processAntiDepsMemory(ins, memory)
	}
}

func processAntiDepsReg(ins *instruction, regs keyInsMap) {
	for r := range ins.outRegs {
		regs[r] = ins
	}

	for r := range ins.inRegs {
		// The instruction might consume the same register we writes. In
		// such a case we don't want the instruction to be anti
		// dependent on itself. If there is any latter instruction
		// writing the same register, such a dependency is not
		// anti-dependency but output dependency.
		dep, ok := regs[r]
		if !ok || dep == ins {
			continue
		}

		addDep(ins, dep)
	}
}

func processAntiDepsMemory(ins *instruction, memory keyInsMap) {
	for _, s := range ins.stores {
		memory[s.Key()] = ins
	}

	for _, l := range ins.loads {
		// The instruction might read memory before it writes it. In
		// such a case we don't want the instruction to be anti
		// dependent on itself. If there is any latter instruction
		// writing the same memory, such a dependency is not anti
		// dependency but output dependency.
		dep, ok := memory[l.Key()]
		if !ok || dep == ins {
			continue
		}

		addDep(ins, dep)
	}
}
