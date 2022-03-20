package riscv

import (
	"decomp/pkg/model"
	"fmt"
	"strings"
)

type Instruction struct {
	value  uint32
	opcode *instructionOpcode
}

func newInstruction(b []byte, opcode *instructionOpcode) Instruction {
	if l := len(b); l < instructionLen {
		panic(fmt.Sprintf("not enough bytes to represent valid opcode: %d", l))
	}

	var value uint32
	for i, v := range b[:instructionLen] {
		value |= uint32(v) << (8 * i)
	}

	return Instruction{
		value:  value,
		opcode: opcode,
	}
}

func (i Instruction) inputRegs() map[model.Register]struct{} {
	if i.opcode.inputRegCnt == 0 {
		return nil
	}

	regs := make(map[model.Register]struct{}, i.opcode.inputRegCnt)

	regs[model.Register(rs1.regNum(i.value))] = struct{}{}
	if i.opcode.inputRegCnt > 1 {
		regs[model.Register(rs2.regNum(i.value))] = struct{}{}
	}

	return regs
}

func (i Instruction) outputRegs() map[model.Register]struct{} {
	if !i.opcode.hasOutputReg {
		return nil
	}

	return map[model.Register]struct{}{
		model.Register(rd.regNum(i.value)): {},
	}
}

func (i Instruction) Value() uint32 { return i.value }
func (i Instruction) Name() string  { return i.opcode.name }

func (i Instruction) String() string {
	// Optimistic preallocation.
	arguments := make([]string, 0, 3)

	if i.opcode.hasOutputReg {
		arguments = append(arguments, rd.regNum(i.value).String())
	}
	if i.opcode.inputRegCnt > 0 {
		arguments = append(arguments, rs1.regNum(i.value).String())
	}
	if i.opcode.inputRegCnt > 1 {
		arguments = append(arguments, rs2.regNum(i.value).String())
	}

	if imm, ok := i.opcode.immediate.parseValue(i.value); ok {
		arguments = append(arguments, fmt.Sprintf("%d", imm))
	}

	return fmt.Sprintf("%s %s", i.opcode.name, strings.Join(arguments, ", "))
}
