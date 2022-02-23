package opcode

import (
	"fmt"
	"sort"
)

type opcode struct {
	getter OpcodeGetter
	opcode Opcode
	masked []byte
}

func duplicateOpcodeErr(op1 opcode, op2 opcode) error {
	return fmt.Errorf("ambiguous opcodes: %s (%s) and %s (%s)",
		op1.getter.Name(), op1.opcode.String(),
		op2.getter.Name(), op2.opcode.String())
}

type Decoder struct {
	groups []maskGroup
}

func NewDecoder(opcs ...OpcodeGetter) (*Decoder, error) {
	opcodes := make([]opcode, len(opcs))
	for i, g := range opcs {
		o := g.Opcode()
		if err := o.Validate(); err != nil {
			return nil, fmt.Errorf("invalid opcode definition %d/%d: %w",
				i, len(opcs), err)
		}

		opcodes[i] = opcode{
			getter: g,
			opcode: o,
			masked: applyMask(o.Bytes, o.Mask),
		}
	}

	groups, err := group(opcodes)
	if err != nil {
		return nil, fmt.Errorf("opcode grouping failed: %w", err)
	}

	return &Decoder{
		groups: groups,
	}, nil
}

func group(opcodes []opcode) ([]maskGroup, error) {
	sort.Slice(opcodes, func(i, j int) bool {
		return byteLT(opcodes[i].opcode.Mask, opcodes[j].opcode.Mask)
	})

	groups := make([]maskGroup, 0, 1)
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

// Match matches a sequence of bytes to Opcode returned by OpcodeGetters. This
// method returns nil if no opcode was matched.
func (d *Decoder) Match(bytes []byte) OpcodeGetter {
	for _, g := range d.groups {
		ins, ok := g.matchInstruction(bytes)
		if ok {
			return ins.getter
		}
	}

	return nil
}
