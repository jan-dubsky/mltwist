package riscv

import (
	"decomp/internal/opcode"
	"decomp/pkg/model"
	"fmt"
)

type ParsingStrategy struct{}

func (*ParsingStrategy) Parse(bytes []byte) (model.Instruction, error) {
	if l := len(bytes); l < instructionLen {
		return model.Instruction{}, fmt.Errorf(
			"bytes are too short (%d) to represent an instruction opcode", l)
	}

	found := decoder64.Match(bytes)
	if found == nil {
		return model.Instruction{}, fmt.Errorf(
			"unknown instruction opcode: 0x%x", bytes[:instructionLen])
	}

	opcode := found.(*instructionOpcode)
	instr := newInstruction(bytes, opcode)

	return model.Instruction{
		Type:    opcode.instrType,
		ByteLen: instructionLen,

		InputRegistry:  instr.inputRegs(),
		OutputRegistry: instr.outputRegs(),

		Details: instr,
	}, nil
}

func newDecoder(opcs ...*instructionOpcode) *opcode.Decoder {
	getters := make([]opcode.OpcodeGetter, len(opcs))
	for i, ins := range opcs {
		getters[i] = ins
	}

	dec, err := opcode.NewDecoder(getters...)
	if err != nil {
		panic(fmt.Sprintf("decoder initialization failed: %s", err.Error()))
	}

	return dec
}

var (
	decoder32 = newDecoder(known32...)
	decoder64 = newDecoder(known64...)
)
