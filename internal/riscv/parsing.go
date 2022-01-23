package riscv

import (
	"decomp/internal/instruction"
	"decomp/internal/opcode"
	"fmt"
)

const instrLen = 4

type ParsingStrategy struct{}

func (*ParsingStrategy) Window() uint64 { return instrLen }

func (*ParsingStrategy) Parse(bytes []byte) (instruction.Instruction, error) {
	found := decoder.Match(bytes)
	if found == nil {
		err := fmt.Errorf("unknown instruction opcode: 0x%x", bytes)
		return instruction.Instruction{}, err
	}

	return instruction.Instruction{
		ByteLen: instrLen,
		Type:    instruction.TypeAritm,
		Details: *found.(*instructionOpcode),
	}, nil
}

var decoder = func() *opcode.Decoder {
	getters := make([]opcode.OpcodeGetter, len(known32))
	for i, ins := range known32 {
		getters[i] = ins
	}

	dec, err := opcode.NewDecoder(getters...)
	if err != nil {
		panic(fmt.Sprintf("decoder initialization failed: %s", err.Error()))
	}

	return dec
}()
