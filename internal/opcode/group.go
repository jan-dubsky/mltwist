package opcode

import (
	"sort"
)

// maskGroup groups all opcodes with the same mask. Mask equality is defined
// by byteEQ function.
type maskGroup[T Opcoder] struct {
	// mask is bit mask shared by all opcodes in the maskGroup.
	mask []byte

	// opcodes is a non-conflicting list of opcodes in the group.
	//
	// The list is sorted by masked opcodes of instruction. This property
	// allows to find an instruction in the group in logarithmic time rather
	// than in linear time. The sort order is defined by byteTL function.
	opcodes []opcode[T]
}

// newMaskGroup creates a new maskGroup consisting of opcs. All opcs are assumed
// to share the same opcode mask, otherwise the behaviour is undefined.
func newMaskGroup[T Opcoder](opcs []opcode[T]) (maskGroup[T], error) {
	sort.Slice(opcs, func(i, j int) bool {
		return byteLT(opcs[i].masked, opcs[j].masked)
	})

	// All masks are the same, so all masked opcodes have the same length.
	// Consequently one opcode cannot be prefix of another and linear check
	// in sorted array is sufficient to guarantee no duplicates.
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

// matchInstruction finds an opcode in the group.
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
