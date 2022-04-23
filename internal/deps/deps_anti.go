package deps

type antiDepProcessor struct {
	regs map[string]*instruction

	memory *instruction
}

func processAntiDeps(instrs []*instruction) {
	p := antiDepProcessor{
		regs: make(map[string]*instruction, numRegs),
	}

	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		p.processRegDeps(ins)
		p.processMemDeps(ins)
	}
}

func (p *antiDepProcessor) processRegDeps(ins *instruction) {
	for r := range ins.outRegs {
		p.regs[r] = ins
	}

	for r := range ins.inRegs {
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
	if len(ins.stores) > 0 {
		p.memory = ins
	}

	if p.memory == nil {
		return
	}

	// The instruction might read memory before it writes it. In such a case
	// we don't want the instruction to be dependent on itself.
	if (len(ins.loads) > 0) && p.memory != ins {
		p.link(ins, p.memory)
	}
}

func (*antiDepProcessor) link(first, second *instruction) {
	first.antiDepsFwd[second] = struct{}{}
	second.antiDepsBack[first] = struct{}{}
}
