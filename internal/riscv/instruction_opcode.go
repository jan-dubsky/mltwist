package riscv

import (
	"decomp/internal/instruction"
	"decomp/internal/opcode"
	"fmt"
)

// instructionLen is length of RISC V opcode in bytes.
const instructionLen = 4

func assertOpcodeLen(b []byte) {
	if l := len(b); l != instructionLen {
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

	instrType instruction.Type
}

func (i instructionOpcode) Opcode() opcode.Opcode { return i.opcode }

func (i instructionOpcode) String() string {
	// FIXME: This is a hack!
	return i.name
}
