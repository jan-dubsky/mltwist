package opcode

import (
	"sort"
)

type maskGroup struct {
	mask    []byte
	opcodes []opcode
}

func newMaskGroup(opcs []opcode) (maskGroup, error) {
	sort.Slice(opcs, func(i, j int) bool {
		return byteLT(opcs[i].masked, opcs[j].masked)
	})

	for i := range opcs[1:] {
		if byteEQ(opcs[i].masked, opcs[i+1].masked) {
			return maskGroup{}, duplicateOpcodeErr(opcs[i], opcs[i+1])
		}
	}

	return maskGroup{
		mask:    opcs[0].opcode.Mask,
		opcodes: opcs,
	}, nil
}

func (g *maskGroup) matchInstruction(bytes []byte) (opcode, bool) {
	if len(g.mask) > len(bytes) {
		return opcode{}, false
	}

	masked := applyMask(bytes[:len(g.mask)], g.mask)
	idx := sort.Search(len(g.opcodes), func(i int) bool {
		curr := g.opcodes[i].masked
		return byteLT(masked, curr) || byteEQ(masked, curr)
	})

	if idx == len(g.opcodes) {
		return opcode{}, false
	}

	opc := g.opcodes[idx]
	if !byteEQ(masked, opc.masked) {
		return opcode{}, false
	}

	return opc, true
}
