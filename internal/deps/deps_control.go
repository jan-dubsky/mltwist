package deps

func processControlDeps(instrs []*instruction) {
	last := instrs[len(instrs)-1]
	for _, ins := range instrs[:len(instrs)-1] {
		ins.controlDepsFwd[last] = struct{}{}
		last.controlDepsBack[ins] = struct{}{}
	}
}
