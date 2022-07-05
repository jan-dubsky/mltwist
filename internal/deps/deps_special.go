package deps

// isMemAccess checks if ins either loads or stores value from or to memory.
func isMemAccess(ins *instruction) bool {
	return len(ins.stores) > 0 || len(ins.loads) > 0
}

// insMemOrder identifies whether an instruction is a memory order special
// instruction type.
func insMemOrder(ins *instruction) bool {
	return ins.typ.MemOrder()
}

// insSpecial identifies if an instruction is a special type of an instruction
// other than memory order.
func insSpecial(ins *instruction) bool {
	t := ins.typ
	return t.Syscall() || t.CPUStateChange()
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
	var (
		lastMemOrder *instruction
		lastSpecial  *instruction
	)

	for _, ins := range instrs {
		if lastMemOrder != nil && isMemAccess(ins) {
			addDep(lastMemOrder, ins)
		}
		if lastSpecial != nil {
			addDep(lastSpecial, ins)
		}

		if insMemOrder(ins) {
			lastMemOrder = ins
		}
		if insSpecial(ins) {
			// This is jump optimization of graph size. As special
			// instruction is dependent on everything we don't need
			// to add dependencies for memory order instructions as
			// well.
			lastMemOrder = nil
			lastSpecial = ins
		}
	}

	// There is no special instruction in the block so the other walk will
	// be no-op. This is just performance optimizations.
	if lastMemOrder == nil && lastSpecial == nil {
		return
	}

	lastMemOrder, lastSpecial = nil, nil
	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]

		if lastMemOrder != nil && (isMemAccess(ins) || insMemOrder(ins)) {
			addDep(ins, lastMemOrder)
		}
		if lastSpecial != nil {
			addDep(ins, lastSpecial)
		}

		if insMemOrder(ins) {
			lastMemOrder = ins
		}
		if insSpecial(ins) {
			// This is jump optimization of graph size. As special
			// instruction is dependent on everything we don't need
			// to add dependencies for memory order instructions as
			// well.
			lastMemOrder = nil
			lastSpecial = ins
		}
	}
}
