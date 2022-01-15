package opcode

import (
	"fmt"
	"sort"
)

type instr struct {
	getter OpcodeGetter
	opcode Opcode
	masked []byte
}

type maskGroup struct {
	mask         []byte
	instructions []instr
}

type Decoder struct {
	groups []maskGroup
}

func NewDecoder(opcodes ...OpcodeGetter) (*Decoder, error) {
	instrs := make([]instr, len(opcodes))
	for i, g := range opcodes {
		o := g.Opcode()
		if err := o.Validate(); err != nil {
			return nil, fmt.Errorf("invalid opcode definition %d/%d: %w",
				i, len(opcodes), err)
		}

		instrs[i] = instr{
			getter: g,
			opcode: o,
			masked: applyMask(o.Bytes, o.Mask),
		}
	}

	return &Decoder{
		groups: group(instrs),
	}, nil
}

func group(instrs []instr) []maskGroup {
	sort.Slice(instrs, func(i, j int) bool {
		return byteLE(instrs[i].opcode.Mask, instrs[j].opcode.Mask)
	})

	groups := make([]maskGroup, 0, uniqueMasks(instrs))
	for begin := 0; begin < len(instrs); {
		mask := instrs[begin].opcode.Mask
		end := begin + sort.Search(len(instrs)-begin, func(i int) bool {
			return !byteEQ(instrs[i+begin].opcode.Mask, mask)
		})

		g := maskGroup{
			mask:         mask,
			instructions: instrs[begin:end],
		}

		groups = append(groups, g)
		begin = end
	}

	return groups
}

func uniqueMasks(instrs []instr) int {
	var (
		last []byte
		cnt  int
	)

	for _, ins := range instrs {
		if !byteEQ(last, ins.opcode.Mask) {
			cnt++
			last = ins.opcode.Mask
		}
	}

	return cnt
}

func (d *Decoder) Match(bytes []byte) OpcodeGetter {
	for _, g := range d.groups {
		ins, ok := g.matchInstruction(bytes)
		if ok {
			return ins.getter
		}
	}

	return nil
}

func (g *maskGroup) matchInstruction(bytes []byte) (instr, bool) {
	if len(g.mask) > len(bytes) {
		return instr{}, false
	}

	masked := applyMask(bytes, g.mask)
	idx := sort.Search(len(g.instructions), func(i int) bool {
		return byteLE(masked, g.instructions[i].masked)
	})

	if idx == len(g.instructions) {
		return instr{}, false
	}
	return g.instructions[idx], true
}
