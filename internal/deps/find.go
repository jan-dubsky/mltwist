package deps

import "decomp/internal/repr"

func Find(seq []repr.Instruction) []*Instruction {
	m := newCPUModel()

	instrs := make([]*Instruction, len(seq))
	for i, ins := range seq {
		instr := newInstruction(ins)
		instr.deps = m.process(instr)

		instrs[i] = instr
	}

	return instrs
}
