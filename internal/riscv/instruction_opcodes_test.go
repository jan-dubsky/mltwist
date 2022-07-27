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

func TestOpcodes32(t *testing.T) {
	assertValidOpcode(t, known32...)
	assertUniqueOpcode(t, known32...)
}

func TestOpcodes64(t *testing.T) {
	assertValidOpcode(t, known64...)
	assertUniqueOpcode(t, known64...)
}

func instructionNameMap(instrs []*instructionOpcode) map[string]*instructionOpcode {
	m := make(map[string]*instructionOpcode, len(instrs))
	for i, ins := range instrs {
		if _, ok := m[ins.name]; ok {
			panic(fmt.Sprintf("duplicate instruction at %d/%d: %s",
				i, len(instrs), ins.name))
		}
		m[ins.name] = ins
	}

	return m
}

var (
	instrMap32 = instructionNameMap(known32)
	instrMap64 = instructionNameMap(known64)
)
