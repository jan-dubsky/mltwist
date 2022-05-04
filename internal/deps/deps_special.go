package deps

func insSpecial(ins *instruction) bool {
	t := ins.typ
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
			addDep(last, ins)
		}

		if specials[i] {
			last = ins
		}
	}

	last = nil
	for i := len(instrs) - 1; i >= 0; i-- {
		ins := instrs[i]
		if last != nil {
			addDep(ins, last)
		}

		if specials[i] {
			last = ins
		}
	}
}
