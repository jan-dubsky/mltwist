package riscv

import (
	"decomp/pkg/model"
	"fmt"
	"strings"
)

type InstrBytes [instructionLen]byte

func (bytes InstrBytes) regNum(r reg) regNum {
	b := bytes[:]
	bitOffsetRaw := r.bitOffset()
	byteOffset := bitOffsetRaw / 8
	b = b[byteOffset:]
	bitOffset := bitOffsetRaw % 8

	var num uint8
	num = b[0] >> bitOffset

	// Bits overlowing to b[1]
	var remBits uint8
	if end := bitOffset + regBits; end > 8 {
		remBits = end - 8
	}

	remVal := b[1] & ((uint8(1) << remBits) - 1)
	num |= remVal << (8 - bitOffset)

	if num >= regCnt {
		panic(fmt.Sprintf("BUG: parsed RISC-V register number: %d", num))
	}

	return regNum(num)
}

func (b InstrBytes) uint32() uint32 {
	var value uint32
	for i, v := range b {
		value |= uint32(v) << (8 * i)
	}

	return value
}

type Instruction struct {
	bytes  InstrBytes
	opcode *instructionOpcode
}

func newInstruction(b []byte, opcode *instructionOpcode) Instruction {
	if l := len(b); l < instructionLen {
		panic(fmt.Sprintf("not enough bytes to represent valid opcode: %d", l))
	}

	var bytes InstrBytes
	copy(bytes[:], b)

	return Instruction{
		bytes:  bytes,
		opcode: opcode,
	}
}

func (i Instruction) inputRegs() map[model.Register]struct{} {
	if i.opcode.inputRegCnt == 0 {
		return nil
	}

	regs := make(map[model.Register]struct{}, i.opcode.inputRegCnt)

	regs[model.Register(i.bytes.regNum(rs1))] = struct{}{}
	if i.opcode.inputRegCnt > 1 {
		regs[model.Register(i.bytes.regNum(rs2))] = struct{}{}
	}

	return regs
}

func (i Instruction) outputRegs() map[model.Register]struct{} {
	if !i.opcode.hasOutputReg {
		return nil
	}

	return map[model.Register]struct{}{
		model.Register(i.bytes.regNum(rd)): {},
	}
}

func (i Instruction) Bytes() InstrBytes { return i.bytes }
func (i Instruction) Name() string      { return i.opcode.name }

func (i Instruction) String() string {
	// Optimistic preallocation.
	arguments := make([]string, 0, 3)

	if i.opcode.hasOutputReg {
		arguments = append(arguments, i.bytes.regNum(rd).String())
	}
	if i.opcode.inputRegCnt > 0 {
		arguments = append(arguments, i.bytes.regNum(rs1).String())
	}
	if i.opcode.inputRegCnt > 1 {
		arguments = append(arguments, i.bytes.regNum(rs2).String())
	}

	if imm, ok := i.opcode.immediate.parseValue(i.bytes); ok {
		arguments = append(arguments, fmt.Sprintf("%d", imm))
	}

	return fmt.Sprintf("%s %s", i.opcode.name, strings.Join(arguments, ", "))
}
