package riscv

import (
	"fmt"
	"mltwist/pkg/model"
	"strings"
)

// instruction represents a parser RISC-V instruction. Unlike the
// instructionType which describes only instruction opcode and properties of the
// opcode, instruction represents the whole instruction including encoding of
// registers, immediate value and last but not least the position in the code.
type instruction struct {
	// addr is virtual address of the instruction in a program memory.
	addr model.Addr

	// value represents the value of the instruction in the memory.
	//
	// As encoding of RISC-V instruction is little-endian, the value encoded
	// is also little endian. In other words, the byte of instruction with
	// the smallest in-memory address is represented by bits [0..7] of the
	// value.
	//
	// As all instructions in RISC-V (excluding compressed instructions
	// which we don't support for the time being) have 4 bytes, we can use
	// uint32 to represent them. Usage of uint32 over []byte has many
	// reasons starting with convenience, followed by performance and ending
	// in memory space consumption.
	value uint32

	// instrType refers the type of the instruction.
	instrType *instructionType
}

// newInstruction crates a new instance of instruction. The new instruction is
// ata address a, is represented of first instructionLen bytes of b and has type
// t.
func newInstruction(a model.Addr, b []byte, t *instructionType) instruction {
	if l := len(b); l < instructionLen {
		panic(fmt.Sprintf("not enough bytes to represent valid opcode: %d", l))
	}

	var value uint32
	for i, v := range b[:instructionLen] {
		value |= uint32(v) << (8 * i)
	}

	return instruction{
		addr:      a,
		value:     value,
		instrType: t,
	}
}

// Name returns name of the instruction.
func (i instruction) Name() string { return i.instrType.name }

// String returns a string representation of an instruction which corresponds to
// standard RISC-V assembler notation of instructions.
func (i instruction) String() string {
	// The value of 3 is Optimistic preallocation.
	as := make([]string, 0, 3)

	if i.instrType.hasOutputReg {
		as = append(as, rd.regNum(i.value).String())
	}
	if i.instrType.inputRegCnt > 0 {
		as = append(as, rs1.regNum(i.value).String())
	}
	if i.instrType.inputRegCnt > 1 {
		as = append(as, rs2.regNum(i.value).String())
	}

	if imm, ok := i.instrType.immediate.parseValue(i.value); ok {
		immStr := fmt.Sprintf("%d", imm)

		// Store instruction is written in order: `s[bhwd] <reg>
		// <mem_addr>`, even though memory base register is r1 and the
		// register written to memory is r2. Why not to make the
		// assembler irregular such that store is the only instruction
		// where direction of data flow is from left operand to the
		// right one. Long live irregularities! Intel celebrates, the
		// rest of the world cries...
		if i.instrType.storeBytes > 0 {
			l := len(as)
			as[l-1], as[l-2] = as[l-2], as[l-1]
		}

		// For some weird reason load and store instructions use
		// different syntax then all other instructions - memory offset
		// syntax.
		if i.instrType.loadBytes > 0 || i.instrType.storeBytes > 0 {
			as[len(as)-1] = fmt.Sprintf("%s(%s)", immStr, as[len(as)-1])
		} else {
			as = append(as, immStr)
		}
	}

	return fmt.Sprintf("%s %s", i.instrType.name, strings.Join(as, ", "))
}
