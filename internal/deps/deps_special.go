package deps

// insSpecial identifies if an instruction is a special type of an instructions.
func insSpecial(ins *instruction) bool {
	t := ins.typ
	return t.Syscall() || t.MemOrder() || t.CPUStateChange()
}

// parseSpecialDeps finds dependencies in between instructions we don't fully
// understand.
//
// Even thought we have description of most of actions of all the instructions,
// there are still instructions which very special meaning which we cannot
// analyze. A good example of such an instruction are atomic instructions which
// might cause either acquire of release semantic synchronization (or both).
// Without any loss of genericity, let's assume that the atomic instruction has
// acquire semantics. In such a case no following memory reads and writes can be
// reordered before this instruction. But as we don't understand acquire
// semantics, we have to be conservative and prohibit any reordering with any
// memory order instruction. Another good example is syscall which in our
// representation has no dependencies. But the OS has the capability to change
// an arbitrary memory address and value of an arbitrary register, so we have to
// prohibit any reordering with the syscall instruction.
func findSpecialDeps(instrs []*instruction) {
	var last *instruction
	for _, ins := range instrs {
		if last != nil {
			addDep(last, ins)
		}

		if insSpecial(ins) {
			last = ins
		}
	}

	last = nil
	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		if last != nil {
			addDep(ins, last)
		}

		if insSpecial(ins) {
			last = ins
		}
	}
}
