package deps

import "mltwist/pkg/expr"

// keyInsMap maps expr.Key to an instruction pointer.
type keyInsMap map[expr.Key]*instruction

// findAntiDeps finds anti dependencies in the code.
//
// Anti dependency is dependency of instruction caused by reusing the same
// register or memory place for completely unrelated value. If one instruction
// consumes a value from register or memory place and later instruction in the
// code writes the same register or memory location, those two instructions are
// dependent on one another. This dependency is a real dependency caused by
// transfer of data, but simply by the fact that the data has to be consumed
// because some other instruction will rewrite them.
func findAntiDeps(instrs []*instruction) {
	regs := make(keyInsMap, numRegs)
	memory := make(keyInsMap, 1)

	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		findAntiDepsReg(ins, regs)
		findAntiDepsMemory(ins, memory)
	}
}

// findAntiDepsRegs finds all register-based anti dependencies in the code.
func findAntiDepsReg(ins *instruction, regs keyInsMap) {
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

// findAntiDepsMemory finds memory-based anti dependencies in the code.
//
// Unlike register data flows where the analysis is trivial, the analysis of
// memory dependencies is more complicated. Namely we'd have to be able to
// evaluate memory address of a particular memory store or load. Unfortunately
// this is not possible during static analysis in a general case the memory
// address might be runtime value. So the best we can do is to introduce some
// heuristics to say that 2 instruction are certainly dependent on one another,
// but we have to be pessimistic and see dependencies wherever we cannot
// guarantee instructions not to be dependent one on another.
//
// For the time being no heuristics is implemented and all stores to one memory
// key are considered to be anti dependent on a previous load.
func findAntiDepsMemory(ins *instruction, memory keyInsMap) {
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
