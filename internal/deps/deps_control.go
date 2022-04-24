package deps

import (
	"mltwist/internal/parser"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

func controlFlowInstruction(ins parser.Instruction) bool {
	for _, t := range ins.JumpTargets {
		if c, ok := t.(expr.Const); ok {
			addr, ok := expr.ConstUint[model.Addr](c)
			if ok && addr == ins.NextAddr() {
				continue
			}
		}

		return true
	}
	return false
}

func processControlDeps(instrs []*instruction) {
	last := instrs[len(instrs)-1]

	// If last instruction is not a jump, the basic block ends there simply
	// because the following instruction in memory is a jump target of some
	// other jump in the code. In such a case it doesn't matter which
	// instruction will be the last in this basic block as none of them is
	// control flow instruction and every single one will result in
	// increment of instruction pointer -> no control dependency at all.
	if !controlFlowInstruction(last.Instr) {
		return
	}

	for _, ins := range instrs[:len(instrs)-1] {
		ins.depsFwd[last] = struct{}{}
		last.depsBack[ins] = struct{}{}
	}
}
