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

func (i instructionOpcode) Validate() error {
	if i.inputRegCnt > 2 {
		return fmt.Errorf("too many input registers: %d", i.inputRegCnt)
	}

	if err := i.Opcode().Validate(); err != nil {
		return fmt.Errorf("invalid opcode definition: %w", err)
	}

	return nil
}

func assertValid(t testing.TB, instrs ...*instructionOpcode) {
	for _, ins := range instrs {
		require.NoError(t, ins.Validate())
	}
}

func assertUnique(t testing.TB, instrs ...*instructionOpcode) {
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
