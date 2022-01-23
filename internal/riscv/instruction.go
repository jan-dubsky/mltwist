package riscv

import (
	"decomp/internal/opcode"
	"fmt"
)

// opcodeLen is length of RISC V opcode in bytes.
const opcodeLen = 4

func assertOpcodeLen(b []byte) {
	if l := len(b); l != opcodeLen {
		panic(fmt.Sprintf("invalid byte slice length to represent an opcode: %d",
			l))
	}
}

type instructionOpcode struct {
	name         string
	opcode       opcode.Opcode
	inputRegCnt  uint8
	hasOutputReg bool
	loadBytes    uint8
	storeBytes   uint8
	unsigned     bool

	immediate immType
}

func (i instructionOpcode) Name() string           { return i.name }
func (i *instructionOpcode) Opcode() opcode.Opcode { return i.opcode }

func (i instructionOpcode) String() string {
	// FIXME: This is a hack!
	return i.Name()
}
