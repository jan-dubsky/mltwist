package opcode

import (
	"fmt"
	"sort"
)

type opcode[T Opcoder] struct {
	getter T
	opcode Opcode
	masked []byte
}

func duplicateOpcodeErr[T Opcoder](op1 opcode[T], op2 opcode[T]) error {
	return fmt.Errorf("ambiguous opcodes: %s (%s) and %s (%s)",
		op1.getter.Name(), op1.opcode.String(),
		op2.getter.Name(), op2.opcode.String())
}

type Decoder[T Opcoder] struct {
	groups []maskGroup[T]
}

func NewDecoder[T Opcoder](opcs ...T) (*Decoder[T], error) {
	opcodes := make([]opcode[T], len(opcs))
	for i, g := range opcs {
		o := g.Opcode()
		if err := o.Validate(); err != nil {
			return nil, fmt.Errorf("invalid opcode definition %d/%d: %w",
				i, len(opcs), err)
		}

		opcodes[i] = opcode[T]{
			getter: g,
			opcode: o,
			masked: applyMask(o.Bytes, o.Mask),
		}
	}

	groups, err := group(opcodes)
	if err != nil {
		return nil, fmt.Errorf("opcode grouping failed: %w", err)
	}

	return &Decoder[T]{
		groups: groups,
	}, nil
}

func group[T Opcoder](opcodes []opcode[T]) ([]maskGroup[T], error) {
	sort.Slice(opcodes, func(i, j int) bool {
		return byteLT(opcodes[i].opcode.Mask, opcodes[j].opcode.Mask)
	})

	groups := make([]maskGroup[T], 0, 1)
	for begin := 0; begin < len(opcodes); {
		mask := opcodes[begin].opcode.Mask
		end := begin + sort.Search(len(opcodes)-begin, func(i int) bool {
			return !byteEQ(opcodes[i+begin].opcode.Mask, mask)
		})

		g, err := newMaskGroup(opcodes[begin:end])
		if err != nil {
			return nil, fmt.Errorf("cannot create group for mask 0x%x: %w",
				mask, err)
		}

		groups = append(groups, g)
		begin = end
	}

	for i, gi := range groups {
		for j, gj := range groups {
			if i == j {
				continue
			}

			for _, o := range gj.opcodes {
				opc, ok := gi.matchInstruction(o.opcode.Bytes)
				if ok {
					return nil, duplicateOpcodeErr(o, opc)
				}
			}
		}
	}

	return groups, nil
}

// Match matches a sequence of bytes to Opcode returned by OpcodeGetters.
func (d *Decoder[T]) Match(bytes []byte) (T, bool) {
	for _, g := range d.groups {
		ins, ok := g.matchInstruction(bytes)
		if ok {
			return ins.getter, true
		}
	}

	var t T
	return t, false
}
