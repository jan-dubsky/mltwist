package riscv

import (
	"decomp/internal/opcode"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type hashableOpcode struct {
	Bytes string
	Mask  string
}

func newHashableOpcode(o opcode.Opcode) hashableOpcode {
	return hashableOpcode{
		Bytes: fmt.Sprintf("%x", o.Bytes),
		Mask:  fmt.Sprintf("%x", o.Mask),
	}
}

func assertValidOpcode(t testing.TB, xlenBytes uint8, instrs []*instructionOpcode) {
	for _, ins := range instrs {
		require.NoError(t, ins.validate(xlenBytes))
	}
}

func assertUniqueOpcode(t testing.TB, instrs []*instructionOpcode) {
	opcodeSet := make(map[hashableOpcode]struct{}, len(instrs))
	for i, ins := range instrs {
		o := ins.Opcode()
		h := newHashableOpcode(o)
		if _, ok := opcodeSet[h]; ok {
			require.Failf(t, "opcode is not unique", "%d: %#v (%#v)",
				i, o, ins)
		}
		opcodeSet[h] = struct{}{}
	}
}

func assertUniqueNames(t testing.TB, instrs []*instructionOpcode) {
	m := make(map[string]*instructionOpcode, len(instrs))
	for _, ins := range instrs {
		_, ok := m[ins.name]
		require.False(t, ok)
		m[ins.name] = ins
	}
}

func TestInstructions(t *testing.T) {
	for arch, exts := range instructions {
		exts := exts
		t.Run(fmt.Sprintf("arch_%v", arch), func(t *testing.T) {
			allInstrs := make([]*instructionOpcode, 0)
			for _, instrs := range exts {
				allInstrs = append(allInstrs, instrs...)
			}

			assertUniqueOpcode(t, allInstrs)
			assertUniqueNames(t, allInstrs)

			for ext, instrs := range exts {
				instrs := instrs
				t.Run(fmt.Sprintf("ext_%v", ext), func(t *testing.T) {
					xlenBytes := 4
					if arch == Variant64 {
						xlenBytes = 8
					}

					assertValidOpcode(t, uint8(xlenBytes), instrs)
				})
			}
		})
	}
}
