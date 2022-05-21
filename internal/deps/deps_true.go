package deps

import "mltwist/pkg/expr"

// numRegs is expected number of registers in a platform.
//
// The purpose of this value is to allow optimistic pre-allocation of maps and
// arrays for the processing. This value doesn't have to be precise neither is
// can be. There is no scientific reasoning or benchmark behind this value, but
// it can serve as a performance optimization.
const numRegs = 32

// findTrueDeps finds true dependencies in the code.
//
// True dependency in between instructions is a form of dependency when one
// instruction consumes output of another instruction. It's obvious that those 2
// instructions are dependent on one another. Data can be passed in between
// instruction using either registers or memory.
func findTrueDeps(instrs []*instruction) {
	regs := make(map[expr.Key]*instruction, numRegs)
	memory := make(map[expr.Key]*instruction, 1)

	for _, ins := range instrs {
		findTrueDepsReg(ins, regs)
		findTrueDepsMemory(ins, memory)
	}
}

// findTrueDepsReg finds all register-based true dependencies in the code.
func findTrueDepsReg(ins *instruction, regs keyInsMap) {
	for r := range ins.inRegs {
		if dep, ok := regs[r]; ok {
			addDep(dep, ins)
		}
	}

	for r := range ins.outRegs {
		regs[r] = ins
	}
}

// findTrueDepsMemory finds memory-based true dependencies in the code.
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
// For the time being no heuristics is implemented and all reads from one memory
// key are considered to be dependent on a previous store.
func findTrueDepsMemory(ins *instruction, memory keyInsMap) {
	// We only see dependencies in between stores and loads. Dependencies in
	// between stores are anti dependencies or output dependencies, but
	// those are not true dependencies.

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
