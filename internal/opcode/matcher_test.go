package opcode

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type opcodeValue struct {
	o Opcode
}

func newValue(bytes []byte, mask []byte) *opcodeValue {
	return &opcodeValue{o: newOpcode(bytes, mask)}
}
func newOpcode(bytes []byte, mask []byte) Opcode { return Opcode{bytes, mask} }
func (o *opcodeValue) Opcode() Opcode            { return o.o }
func (o opcodeValue) Name() string {
	return fmt.Sprintf("test instruction: %s", o.o.String())
}

var testOpcs = []*opcodeValue{
	newValue([]byte{0x25, 0x50, 0x45, 0x88}, []byte{0xff, 0, 0, 0xff}),
	newValue([]byte{0x26, 0x50, 0x45, 0x88}, []byte{0xff, 0, 0, 0xff}),
	newValue([]byte{0x25, 0x50, 0x45, 0x89}, []byte{0xff, 0, 0, 0xff}),
	newValue([]byte{0xa5, 0x50, 0x45, 0x89}, []byte{0xff, 0, 0xf, 0xff}),
	newValue([]byte{0xa5, 0x50, 0x46, 0x89}, []byte{0xff, 0, 0xf, 0xff}),
	newValue([]byte{0xa5, 0x50, 0x58, 0x89}, []byte{0xff, 0, 0x1f, 0xff}),
	newValue([]byte{0xa5, 0x50, 0x58, 0x88}, []byte{0xff, 0, 0x1f, 0xff}),
	newValue([]byte{0xf5, 0x50, 0x45, 0x88}, []byte{0xf0, 0, 0, 0xff}),
	newValue([]byte{0xf5, 0x50, 0x45, 0x89}, []byte{0xf0, 0, 0, 0xff}),
}

func TestMatcher_New_Successful(t *testing.T) {
	r := require.New(t)

	dec, err := NewMatcher(testOpcs)
	r.NoError(err)
	r.Len(dec.groups, 4)

	r.Equal([]byte{0xf0, 0, 0, 0xff}, dec.groups[0].mask)
	r.Equal([]byte{0xff, 0, 0, 0xff}, dec.groups[1].mask)
	r.Equal([]byte{0xff, 0, 0xf, 0xff}, dec.groups[2].mask)
	r.Equal([]byte{0xff, 0, 0x1f, 0xff}, dec.groups[3].mask)

	r.Len(dec.groups[0].opcodes, 2)
	r.Len(dec.groups[1].opcodes, 3)
	r.Len(dec.groups[2].opcodes, 2)
	r.Len(dec.groups[3].opcodes, 2)

	checkSorted := func(opcs []opcode[*opcodeValue]) {
		sorted := sort.SliceIsSorted(opcs, func(i, j int) bool {
			return byteLT(opcs[i].masked, opcs[j].masked)
		})
		r.True(sorted)
	}
	checkSorted(dec.groups[0].opcodes)
	checkSorted(dec.groups[1].opcodes)
	checkSorted(dec.groups[2].opcodes)
	checkSorted(dec.groups[3].opcodes)
}

func TestMatcher_New_Failed(t *testing.T) {
	tests := []struct {
		name string
		opcs []*opcodeValue
	}{{
		name: "ambiguous_same_mask",
		opcs: []*opcodeValue{
			newValue([]byte{0xf5, 0x50, 0x45, 0x88}, []byte{0xf0, 0, 0, 0xff}),
			newValue([]byte{0xf6, 0x50, 0x45, 0x88}, []byte{0xf0, 0, 0, 0xff}),
		},
	}, {
		name: "ambiguous_different_mask",
		opcs: []*opcodeValue{
			newValue([]byte{0xf5, 0x50, 0x45, 0x88}, []byte{0xf0, 0, 0, 0xff}),
			newValue([]byte{0xf6, 0x50, 0x45, 0x88}, []byte{0xff, 0, 0, 0xff}),
		},
	}, {
		name: "invalid_opcode",
		opcs: []*opcodeValue{
			newValue([]byte{0xf5, 0x50, 0x45, 0x88}, []byte{0xf0, 0, 0, 0}),
			newValue([]byte{0xf6, 0x50, 0x45, 0x88}, []byte{0xff, 0, 0, 0xff}),
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			dec, err := NewMatcher(tt.opcs)
			r.Error(err)
			r.Nil(dec)
		})
	}
}

func TestMatcher_Match(t *testing.T) {
	dec, err := NewMatcher(testOpcs)
	require.NoError(t, err)

	tests := []struct {
		bytes    []byte
		expected *opcodeValue
	}{{
		bytes:    []byte{0x25, 0xff, 0xee, 0x88},
		expected: testOpcs[0],
	}, {
		bytes:    []byte{0x25, 0xff, 0xee, 0x89},
		expected: testOpcs[2],
	}, {
		bytes:    []byte{0xa5, 0xff, 0x85, 0x89},
		expected: testOpcs[3],
	}, {
		bytes:    []byte{0xa5, 0x88, 0x98, 0x89},
		expected: testOpcs[5],
	}, {
		bytes:    []byte{0xf5, 0x50, 0x45, 0x89},
		expected: testOpcs[8],
	}, {
		bytes:    []byte{0xaa, 0xaa},
		expected: nil,
	}}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := require.New(t)

			opc, ok := dec.Match(tt.bytes)
			if tt.expected == nil {
				r.False(ok)
				return
			}

			r.True(ok)
			r.Equal(tt.expected, opc)
		})
	}
}
