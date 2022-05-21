package deps

import "mltwist/pkg/expr"

// findOutputDeps finds output dependencies in the code.
//
// Output dependency is a dependency in between instructions which write the
// same register or memory address. As the final result of a sequence of
// instructions must remain always the same, we cannot switch the latest write
// into a register or into a memory place with previous writes to the same
// memory place or register. Such kind of dependency is output dependency.
func findOutputDeps(instrs []*instruction) {
	regs := make(map[expr.Key]*instruction, numRegs)
	memory := make(map[expr.Key]*instruction, 1)

	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		findOutputDepsReg(ins, regs)
		findOutputDepsMemory(ins, memory)
	}
}

// findOutputDepsReg finds register-based output dependencies in the code.
func findOutputDepsReg(ins *instruction, regs keyInsMap) {
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

// findOutputDepsMemory finds memory-based output dependencies in the code.
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
// For the time being no such heuristic is implemented. Consequently we have to
// consider all stores dependent on one another. Due to different memory
// addresses and store writes, the final state of the memory might depend on all
// stores in the sequence.
func findOutputDepsMemory(ins *instruction, memory keyInsMap) {
	// Please note that being dependent is a transitive relation.
	// Consequently it's sufficient for a store to be dependent on the
	// previous store and transitivity then makes it dependent in all stores
	// before.

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
