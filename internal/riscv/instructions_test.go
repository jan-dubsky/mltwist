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

func assertValid(t testing.TB, instrs ...*instr) {
	for _, ins := range instrs {
		require.NoError(t, ins.Opcode().Validate())
	}
}

func assertUnique(t testing.TB, instrs ...*instr) {
	opcodeSet := make(map[hashableOpcode]struct{}, len(known32))
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

func TestOpcodes32(t *testing.T) {
	assertValid(t, known32...)
	assertUnique(t, known32...)
}

func TestOpcodes64(t *testing.T) {
	assertValid(t, known64...)
	assertUnique(t, known64...)
}
