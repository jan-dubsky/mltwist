package opcode

import (
	"sort"
)

type maskGroup[T Opcoder] struct {
	mask    []byte
	opcodes []opcode[T]
}

func newMaskGroup[T Opcoder](opcs []opcode[T]) (maskGroup[T], error) {
	sort.Slice(opcs, func(i, j int) bool {
		return byteLT(opcs[i].masked, opcs[j].masked)
	})

	for i := range opcs[1:] {
		if byteEQ(opcs[i].masked, opcs[i+1].masked) {
			return maskGroup[T]{}, duplicateOpcodeErr(opcs[i], opcs[i+1])
		}
	}

	return maskGroup[T]{
		mask:    opcs[0].opcode.Mask,
		opcodes: opcs,
	}, nil
}

func (g *maskGroup[T]) matchInstruction(bytes []byte) (opcode[T], bool) {
	if len(g.mask) > len(bytes) {
		return opcode[T]{}, false
	}

	masked := applyMask(bytes[:len(g.mask)], g.mask)
	idx := sort.Search(len(g.opcodes), func(i int) bool {
		curr := g.opcodes[i].masked
		return byteLT(masked, curr) || byteEQ(masked, curr)
	})

	if idx == len(g.opcodes) {
		return opcode[T]{}, false
	}

	opc := g.opcodes[idx]
	if !byteEQ(masked, opc.masked) {
		return opcode[T]{}, false
	}

	return opc, true
}
