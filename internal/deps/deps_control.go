package deps

// findControlDeps finds control dependencies in the code.
//
// Control dependency is a dependency in between an instruction and its
// surrounding jumps. The reason is that those jumps determine whether the
// instruction executes or not. As this analysis is supposed to run for a single
// basic block, all previous jumps are already sorted out - those are in another
// basic blocks. But the basic block might still end with a jump instruction and
// in such a case the jump instruction as to be the last in the block. This
// implies that every instruction in the basic block is dependent on the jump
// instruction at the block end.
func findControlDeps(instrs []*instruction) {
	last := instrs[len(instrs)-1]

	// If last instruction is not a jump, the basic block ends there simply
	// because the following instruction in memory is a jump target of some
	// other jump in the code. In such a case it doesn't matter which
	// instruction will be the last in this basic block as none of them is
	// control flow instruction and every single one will result in
	// increment of instruction pointer -> no control dependency at all.
	if len(last.jumpTargets) == 0 {
		return
	}

	for _, ins := range instrs[:len(instrs)-1] {
		addDep(ins, last)
	}
}
