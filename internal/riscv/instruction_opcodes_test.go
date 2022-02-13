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

func validateInstructionOpcode(ins instructionOpcode) error {
	if ins.inputRegCnt > 2 {
		return fmt.Errorf("too many input registers: %d", ins.inputRegCnt)
	}

	if err := ins.Opcode().Validate(); err != nil {
		return fmt.Errorf("invalid opcode definition: %w", err)
	}

	return nil
}

func assertValidOpcode(t testing.TB, instrs ...*instructionOpcode) {
	for _, ins := range instrs {
		require.NoError(t, validateInstructionOpcode(*ins))
	}
}

func assertUniqueOpcode(t testing.TB, instrs ...*instructionOpcode) {
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

func assertUniqueNames(t testing.TB, instrs ...*instructionOpcode) {
	m := make(map[string]*instructionOpcode, len(instrs))
	for _, ins := range instrs {
		_, ok := m[ins.name]
		require.False(t, ok)
		m[ins.name] = ins
	}
}

func TestOpcodes32(t *testing.T) {
	assertValidOpcode(t, known32...)
	assertUniqueOpcode(t, known32...)
	assertUniqueNames(t, known32...)
}

func TestOpcodes64(t *testing.T) {
	assertValidOpcode(t, known64...)
	assertUniqueOpcode(t, known64...)
	assertUniqueNames(t, known64...)
}
