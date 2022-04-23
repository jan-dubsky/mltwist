package deps

// numRegs is expected number of registers in a platform.
//
// The purpose of this value is to allow optimistic pre-allocation of maps and
// arrays for the processing. This value doesn't have to be precise neither is
// can be. There is no scientific reasoning or benchmark behind this value, but
// it can serve as a performance optimization.
const numRegs = 32

type trueDepProcessor struct {
	regs map[string]*instruction

	memory *instruction
}

func processTrueDeps(instrs []*instruction) {
	p := trueDepProcessor{
		regs: make(map[string]*instruction, numRegs),
	}

	for _, ins := range instrs {
		p.processRegDeps(ins)
		p.processMemDeps(ins)
	}
}

func (p *trueDepProcessor) processRegDeps(ins *instruction) {
	for r := range ins.inRegs {
		if dep, ok := p.regs[r]; ok {
			p.link(dep, ins)
		}
	}

	for r := range ins.outRegs {
		p.regs[r] = ins
	}
}

func (p *trueDepProcessor) processMemDeps(ins *instruction) {
	if p.memory != nil && len(ins.loads) > 0 {
		p.link(p.memory, ins)
	}

	if len(ins.stores) > 0 {
		p.memory = ins
	}
}

func (*trueDepProcessor) link(first, second *instruction) {
	first.trueDepsFwd[second] = struct{}{}
	second.trueDepsBack[first] = struct{}{}
}
