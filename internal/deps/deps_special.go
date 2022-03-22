package deps

func insSpecial(ins *instruction) bool {
	t := ins.Instr.Type
	return t.Syscall() || t.MemOrder() || t.CPUStateChange()
}

func processSpecialDeps(instrs []*instruction) {
	specials := make([]bool, len(instrs))
	for i, ins := range instrs {
		specials[i] = insSpecial(ins)
	}

	var last *instruction
	for i, ins := range instrs {
		if last != nil {
			last.specialDepsFwd[ins] = struct{}{}
			ins.specialDepsBack[last] = struct{}{}
		}

		if specials[i] {
			last = ins
		}
	}

	last = nil
	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		if last != nil {
			ins.specialDepsFwd[last] = struct{}{}
			last.specialDepsBack[ins] = struct{}{}
		}

		if specials[i] {
			last = ins
		}
	}
}
