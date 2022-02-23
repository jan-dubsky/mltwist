package riscv

import (
	"decomp/internal/instruction"
	"decomp/internal/opcode"
)

// instructionLen is length of RISC V opcode in bytes.
const instructionLen = 4

type instructionOpcode struct {
	name   string
	opcode opcode.Opcode

	inputRegCnt  uint8
	hasOutputReg bool

	loadBytes  uint8
	storeBytes uint8

	unsigned bool

	immediate           immType
	additionalImmediate addOpcodeInfo

	instrType instruction.Type
}

func (i instructionOpcode) Opcode() opcode.Opcode { return i.opcode }
func (i instructionOpcode) Name() string          { return i.name }
