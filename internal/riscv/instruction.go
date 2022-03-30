package riscv

import (
	"decomp/pkg/model"
	"fmt"
	"strings"
)

type Instruction struct {
	address model.Address
	value   uint32
	opcode  *instructionOpcode
}

func newInstruction(a model.Address, b []byte, opcode *instructionOpcode) Instruction {
	if l := len(b); l < instructionLen {
		panic(fmt.Sprintf("not enough bytes to represent valid opcode: %d", l))
	}

	var value uint32
	for i, v := range b[:instructionLen] {
		value |= uint32(v) << (8 * i)
	}

	return Instruction{
		address: a,
		value:   value,
		opcode:  opcode,
	}
}

func (i Instruction) inputRegs() model.Registers {
	if i.opcode.inputRegCnt == 0 {
		return nil
	}

	regs := make(model.Registers, i.opcode.inputRegCnt)

	regs[model.Register(rs1.regNum(i.value))] = struct{}{}
	if i.opcode.inputRegCnt > 1 {
		regs[model.Register(rs2.regNum(i.value))] = struct{}{}
	}

	return regs
}

func (i Instruction) outputRegs() model.Registers {
	if !i.opcode.hasOutputReg {
		return nil
	}

	return model.Registers{
		model.Register(rd.regNum(i.value)): {},
	}
}

func (i Instruction) name() string { return i.opcode.name }

func (i Instruction) String() string {
	// Optimistic preallocation.
	as := make([]string, 0, 3)

	if i.opcode.hasOutputReg {
		as = append(as, rd.regNum(i.value).String())
	}
	if i.opcode.inputRegCnt > 0 {
		as = append(as, rs1.regNum(i.value).String())
	}
	if i.opcode.inputRegCnt > 1 {
		as = append(as, rs2.regNum(i.value).String())
	}

	if imm, ok := i.opcode.immediate.parseValue(i.value); ok {
		immStr := fmt.Sprintf("%d", imm)

		// Store instruction is written in order: `s[bhwd] <reg>
		// <mem_addr>`, even though memory base register is r1 and the
		// register written to memory is r2. Why not to make the
		// assember irregular such that store is the only instruction
		// where direction of data flow is from left operand to the
		// right one. Long live irregularities! Intel celebrates, the
		// rest of the world cries...
		if i.opcode.storeBytes > 0 {
			l := len(as)
			as[l-1], as[l-2] = as[l-2], as[l-1]
		}

		// For some weird reason load and store instructions use
		// different syntax then all other instructions - memory offset
		// syntax.
		if i.opcode.loadBytes > 0 || i.opcode.storeBytes > 0 {
			as[len(as)-1] = fmt.Sprintf("%s(%s)", immStr, as[len(as)-1])
		} else {
			as = append(as, immStr)
		}
	}

	return fmt.Sprintf("%s %s", i.opcode.name, strings.Join(as, ", "))
}
