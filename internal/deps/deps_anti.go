package deps

import "decomp/pkg/model"

type antiDepProcessor struct {
	regs map[model.Register]*instruction

	memory *instruction
}

func processAntiDeps(instrs []*instruction) {
	p := antiDepProcessor{
		regs: make(map[model.Register]*instruction, numRegs),
	}

	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		p.processRegDeps(ins)
		p.processMemDeps(ins)
	}
}

func (p *antiDepProcessor) processRegDeps(ins *instruction) {
	for r := range ins.Instr.OutputRegistry {
		p.regs[r] = ins
	}

	for r := range ins.Instr.InputRegistry {
		i, ok := p.regs[r]
		if !ok {
			continue
		}

		// The instruction might consume the same register we writes. In
		// such a case we don't want the instruction to be dependent on
		// itself.
		if ins == i {
			continue
		}

		p.link(ins, i)
	}
}

func (p *antiDepProcessor) processMemDeps(ins *instruction) {
	if ins.Instr.Type.Store() {
		p.memory = ins
	}

	if p.memory == nil {
		return
	}

	// The instruction might read memory before it writes it. In such a case
	// we don't want the instruction to be dependent on itself.
	if ins.Instr.Type.Load() && p.memory != ins {
		p.link(ins, p.memory)
	}
}

func (*antiDepProcessor) link(first, second *instruction) {
	first.antiDepsFwd[second] = struct{}{}
	second.antiDepsBack[first] = struct{}{}
}
