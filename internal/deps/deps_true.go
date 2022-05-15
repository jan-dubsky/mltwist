package deps

import "mltwist/pkg/expr"

// numRegs is expected number of registers in a platform.
//
// The purpose of this value is to allow optimistic pre-allocation of maps and
// arrays for the processing. This value doesn't have to be precise neither is
// can be. There is no scientific reasoning or benchmark behind this value, but
// it can serve as a performance optimization.
const numRegs = 32

func processTrueDeps(instrs []*instruction) {
	regs := make(map[expr.Key]*instruction, numRegs)
	memory := make(map[expr.Key]*instruction, 1)

	for _, ins := range instrs {
		processTrueDepsReg(ins, regs)
		processTrueDepsMemory(ins, memory)
	}
}

func processTrueDepsReg(ins *instruction, regs keyInsMap) {
	for r := range ins.inRegs {
		if dep, ok := regs[r]; ok {
			addDep(dep, ins)
		}
	}

	for r := range ins.outRegs {
		regs[r] = ins
	}
}

func processTrueDepsMemory(ins *instruction, memory keyInsMap) {
	for _, l := range ins.loads {
		dep, ok := memory[l.Key()]
		if !ok {
			continue
		}

		addDep(dep, ins)
	}

	for _, s := range ins.stores {
		memory[s.Key()] = ins
	}
}
