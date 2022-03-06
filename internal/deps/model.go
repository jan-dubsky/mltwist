package deps

import "decomp/pkg/model"

type cpuModel struct {
	regs map[model.Register]*Instruction
}

func newCPUModel() cpuModel {
	return cpuModel{
		// 32 registers is quite typical number of registers.
		regs: make(map[model.Register]*Instruction, 32),
	}
}

func (m *cpuModel) process(ins *Instruction) []*Instruction {
	deps := make([]*Instruction, 0)

	for _, r := range ins.Instr.InputRegistry {
		if idx, ok := m.regs[r]; ok {
			deps = append(deps, idx)
		}
	}

	for _, r := range ins.Instr.OutputRegistry {
		m.regs[r] = ins
	}

	return deps
}
