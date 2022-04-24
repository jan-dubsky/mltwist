package deps

type outputDepProcessor struct {
	regs map[string]*instruction

	memory []*instruction
}

func processOutputDeps(instrs []*instruction) {
	p := outputDepProcessor{
		regs: make(map[string]*instruction, numRegs),
	}

	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		p.processRegDeps(ins)
		p.processMemDeps(ins)
	}
}

func (p *outputDepProcessor) processRegDeps(ins *instruction) {
	for r := range ins.outRegs {
		i, ok := p.regs[r]
		if !ok {
			p.regs[r] = ins
			continue
		}

		// We are certain that i != ins.
		p.link(ins, i)
	}
}

func (p *outputDepProcessor) processMemDeps(ins *instruction) {
	if len(ins.stores) == 0 {
		return
	}

	for _, i := range p.memory {
		p.link(ins, i)
	}
	p.memory = append(p.memory, ins)
}

func (*outputDepProcessor) link(first, second *instruction) {
	first.depsFwd[second] = struct{}{}
	second.depsBack[first] = struct{}{}
}
