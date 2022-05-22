package opcode

import (
	"fmt"
	"sort"
)

// opcode is a wrapper structure on top of Opcoder which allows
type opcode[T Opcoder] struct {
	// opcoder is the Opcoder this opcode wraps.
	opcoder T

	// opcode is cached opcode returned by opcoder which speeds up execution
	// of the algorithms - getting the opcode from opcoder has overhead of
	// virtual function call.
	opcode Opcode

	// masked represents opcode bytes anded byte-by-byte with opcode mask.
	//
	// As the definition of Opcode struct allows masked bits to be set in
	// Bytes array in Opcode struct, we have to and mask those bits first to
	// represent real opcode bits.
	masked []byte
}

// duplicateOpcodeErr is a helper function to produce an error stating that op1
// and op2 escribe the same opcode bits - i.e. that those are conflicting
// opcodes.
func duplicateOpcodeErr[T Opcoder](op1 opcode[T], op2 opcode[T]) error {
	return fmt.Errorf("ambiguous opcodes: %s (%s) and %s (%s)",
		op1.opcoder.Name(), op1.opcode.String(),
		op2.opcoder.Name(), op2.opcode.String())
}

// Matcher is an object identifying instruction in the code.
type Matcher[T Opcoder] struct {
	groups []maskGroup[T]
}

// NewMatcher creates new Matcher configured to recognize instruction opcodes in
// opcs.
//
// The opcs array is left untouched and the caller is allowed to further use it.
func NewMatcher[T Opcoder](opcs []T) (*Matcher[T], error) {
	for i, opc := range opcs {
		err := opc.Opcode().Validate()
		if err != nil {
			return nil, fmt.Errorf("invalid opcode definition %d/%d: %w",
				i, len(opcs), err)
		}
	}

	opcodes := newOpcodes(opcs)

	groups, err := group(opcodes)
	if err != nil {
		return nil, fmt.Errorf("opcode grouping failed: %w", err)
	}

	return &Matcher[T]{
		groups: groups,
	}, nil
}

// newOpcodes wraps array of Opcoders into an array of opcodes.
func newOpcodes[T Opcoder](opcs []T) []opcode[T] {
	opcodes := make([]opcode[T], len(opcs))
	for i, opc := range opcs {
		o := opc.Opcode()

		opcodes[i] = opcode[T]{
			opcoder: opc,
			opcode:  o,
			masked:  applyMask(o.Bytes, o.Mask),
		}
	}

	return opcodes
}

// group splits opcodes into maskGroups and validates that no instructions
// conflict one to another.
func group[T Opcoder](opcodes []opcode[T]) ([]maskGroup[T], error) {
	sort.Slice(opcodes, func(i, j int) bool {
		return byteLT(opcodes[i].opcode.Mask, opcodes[j].opcode.Mask)
	})

	groups := make([]maskGroup[T], 0, 1)
	for len(opcodes) > 0 {
		mask := opcodes[0].opcode.Mask
		end := sort.Search(len(opcodes), func(i int) bool {
			return !byteEQ(opcodes[i].opcode.Mask, mask)
		})

		g, err := newMaskGroup(opcodes[:end])
		if err != nil {
			return nil, fmt.Errorf("cannot create group for mask 0x%x: %w",
				mask, err)
		}

		groups = append(groups, g)
		opcodes = opcodes[end:]
	}

	err := checkConflicts(groups)
	if err != nil {
		return nil, fmt.Errorf("instructions in between groups conflict: %w", err)
	}

	return groups, nil
}

// checkConflicts asserts that no instruction conflicts with one another.
//
// The non-conflicting check is to be full n^2 algorithm. Please note that we
// cannot match only for j which is greater than i as instructions can be prefix
// of one another. In other words, the relation of being conflicting is in
// general non-symmetrical. This holds even in case all instructions have the
// same length as mask of one instruction can be bitwise subset of another mask.
func checkConflicts[T Opcoder](groups []maskGroup[T]) error {
	// Make sure that no pair of opcodes conflicts.
	//
	// This has to be full n^2 algorithm - we cannot match only for j which
	// is greater than i as instructions can be prefix of one another. In
	// other words, the relation of being conflicting is in general
	// non-symmetrical.
	for i, gi := range groups {
		for j, gj := range groups {
			if i == j {
				continue
			}

			for _, o := range gj.opcodes {
				opc, ok := gi.matchInstruction(o.opcode.Bytes)
				if ok {
					return duplicateOpcodeErr(o, opc)
				}
			}
		}
	}

	return nil
}

// Match matches a sequence of bytes to an instruction opcode.
//
// It's allowed to pass bs of arbitrary length and those opcodes which can fit
// bs will be matched.
func (d *Matcher[T]) Match(bs []byte) (T, bool) {
	for _, g := range d.groups {
		ins, ok := g.matchInstruction(bs)
		if ok {
			return ins.opcoder, true
		}
	}

	// This is really the smartest way you were able to find to express zero
	// value of a generic parameter? How about some default keyword you have
	// reserved in the language specs?
	var t T
	return t, false
}
