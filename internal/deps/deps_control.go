package deps

import "mltwist/pkg/model"

func controlFlowInstruction(t model.Type) bool {
	return t.Jump() || t.CJump() || t.JumpDyn()
}

func processControlDeps(instrs []*instruction) {
	last := instrs[len(instrs)-1]

	// If last instruction is not a jump, the basic block ends there simply
	// because the following instruction in memory is a jump target of some
	// other jump in the code. In such a case it doesn't matter which
	// instruction will be the last in this basic block as none of them is
	// control flow instruction and every single one will result in
	// increment of instruction pointer -> no control dependency at all.
	if !controlFlowInstruction(last.Instr.Type) {
		return
	}

	for _, ins := range instrs[:len(instrs)-1] {
		ins.controlDepsFwd[last] = struct{}{}
		last.controlDepsBack[ins] = struct{}{}
	}
}
