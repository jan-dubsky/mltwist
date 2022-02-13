package riscv

import (
	"decomp/internal/instruction"
	"decomp/internal/opcode"
	"fmt"
	"sort"
	"strings"
)

const (
	low7Bits byte = 0x7F
	low3Bits byte = 0x7
)

func assertMask(b byte, mask byte) {
	if b&mask != b {
		panic(fmt.Sprintf("bits must match mask 0x%x: 0x%x", mask, b))
	}
}

// revertBytes is an utility function which reverts b as a slice. For more
// convenience, this function also returns b.
//
// The purpose of this method is to lower cognitive complexity of RISC-V
// instruction definition in the code. The problem is that even though RISC-V
// instructions are encoded in little endian byte order, the RICS-V
// specification ifself user big endian notation (least significant byte is the
// right-most byte). This misalignment in between ste specs and real
// implementation is expected. In human readable documents, it's just common to
// write least significant byte and bit to right. But as humans are not good in
// reverting byte sequences they read, it's much better solution to write
// instruction opcodes in big endian and then to revert the array.
func revertBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func opcode7(low byte) opcode.Opcode {
	assertMask(low, low7Bits)

	return opcode.Opcode{
		Bytes: []byte{low},
		Mask:  []byte{low7Bits},
	}
}

func opcode10(high byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(high, low3Bits)

	return opcode.Opcode{
		Bytes: []byte{low, high << 4},
		Mask:  []byte{low7Bits, low3Bits << 4},
	}
}

func opcode17(high byte, mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)
	assertMask(high, low7Bits)

	return opcode.Opcode{
		Bytes: []byte{low, mid << 4, 0, high << 1},
		Mask:  []byte{low7Bits, low3Bits << 4, 0, low7Bits << 1},
	}
}

// opcodeShiftImm creates an opcode definition for RISC V bit shift instruction
// with shift immediate encoded in an instruction opcode.
//
// Even though RISC V manual states that there are only 6 distint instruction
// types and all of them should be describable by either opcode7, opcode10, or
// opcode17, there is one small exception (yes x86, you are not the only one
// architecture doing this...). The exception are bit shift with fixes argument
// in instruction opcode immediate. Technically, such an instruction is just a
// bit shift with 12 bit immediate, but there are a few catches.
//
// The first catch is that not all immediate values are allowed. To be more
// specific, on an architecture with XLEN bits in registers (for simplicity
// let's consider XLEN=32, which is true for 32 bit processors), it doesn't make
// sense to encode more than 32 bit immediate value to shift and though all (but
// one, see below) higher bits of immediate value are reserved to be zero.
//
// Another exception is the logical vs arithmetic shift difference, as the bit
// differentiating those two is bit [30] of an instruction opcode, which
// corresponds to 11th bit of the 12 bit immediate value in I opcode encoding
// type. This brings a weird inconsistency when two instructions represented by
// two different assembler instructions ("srli" and "srai") have the same opcode
// but are differentiated only by a single bit in an immediate value field.
//
// As different immediate value can encode different assembler instructions, we
// need them to be parsed as 2 different instructions. Consequently we are
// forced to describe this instruction opcode meta-format which is not specified
// by the architecture specification document by itself, but which allows us to
// parse bit shifts by an immediate value.
func opcodeShiftImm(arithmetic bool, shiftBits uint8, mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)

	if s := shiftBits; s != 5 && s != 6 {
		panic(fmt.Sprintf("invalid immediate-encoded shift bit count: %d", s))
	}

	var high byte = 0
	if arithmetic {
		high = byte(1) << 6
	}

	// Shift is encoded in bits [20:(20+shiftBits)]. So we do 1<<shiftBits
	// to get 2^shiftBits. Then we subtract 1 which creates is a bit mask
	// for bits encoding values 0..(2^shiftBits)-1. We then invert the mask
	// to force all reserved bits to be zero.
	//
	// And then we have to shift this mask to the right place - to 20th bit
	// of opcode. As we have just high half of instruction opcode, we are
	// already shifted 16 bits. So 4 bits are remaining.
	shiftBitMask := (uint16(1) << shiftBits) - 1
	highHalfMask := (^shiftBitMask) << 4

	return opcode.Opcode{
		Bytes: []byte{high, 0, mid, low},
		Mask: []byte{
			byte(highHalfMask >> 8),
			byte(highHalfMask),
			low3Bits << 4,
			low7Bits,
		},
	}
}

var arithm32 = []*instructionOpcode{
	{
		name:         "lui",
		opcode:       opcode7(0b0110111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "auipc",
		opcode:       opcode7(0b0010111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "jal",
		opcode:       opcode7(0b1101111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeJ,
		instrType:    instruction.TypeJump,
	}, {
		name:         "jalr",
		opcode:       opcode10(0b000, 0b1100111),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeJumpDyn,
	}, {
		name:         "beq",
		opcode:       opcode10(0b000, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    instruction.TypeCJump,
	}, {
		name:         "bne",
		opcode:       opcode10(0b001, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    instruction.TypeCJump,
	}, {
		name:         "blt",
		opcode:       opcode10(0b100, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    instruction.TypeCJump,
	}, {
		name:         "bge",
		opcode:       opcode10(0b101, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    instruction.TypeCJump,
	}, {
		name:         "bltu",
		opcode:       opcode10(0b110, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    instruction.TypeCJump,
	}, {
		name:         "bgeu",
		opcode:       opcode10(0b111, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    instruction.TypeCJump,
	}, {
		name:         "lb",
		opcode:       opcode10(0b000, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		instrType:    instruction.TypeLoad,
	}, {
		name:         "lh",
		opcode:       opcode10(0b001, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		immediate:    immTypeI,
		instrType:    instruction.TypeLoad,
	}, {
		name:         "lw",
		opcode:       opcode10(0b010, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		immediate:    immTypeI,
		instrType:    instruction.TypeLoad,
	}, {
		name:         "lbu",
		opcode:       opcode10(0b100, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    instruction.TypeLoad,
	}, {
		name:         "lhu",
		opcode:       opcode10(0b101, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    instruction.TypeLoad,
	}, {
		name:         "sb",
		opcode:       opcode10(0b000, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   1,
		immediate:    immTypeS,
		instrType:    instruction.TypeStore,
	}, {
		name:         "sh",
		opcode:       opcode10(0b001, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   2,
		immediate:    immTypeS,
		instrType:    instruction.TypeStore,
	}, {
		name:         "sw",
		opcode:       opcode10(0b010, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   4,
		immediate:    immTypeS,
		instrType:    instruction.TypeStore,
	}, {
		name:         "addi",
		opcode:       opcode10(0b000, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "slti",
		opcode:       opcode10(0b010, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "sltiu",
		opcode:       opcode10(0b011, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "xori",
		opcode:       opcode10(0b100, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "ori",
		opcode:       opcode10(0b110, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "andi",
		opcode:       opcode10(0b111, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeAritm,
	}, {
		name:                "slli",
		opcode:              opcodeShiftImm(false, 5, 0b001, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           instruction.TypeAritm,
	}, {
		name:                "srli",
		opcode:              opcodeShiftImm(false, 5, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           instruction.TypeAritm,
	}, {
		name:                "srai",
		opcode:              opcodeShiftImm(true, 5, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           instruction.TypeAritm,
	}, {
		name:         "add",
		opcode:       opcode17(0b0000000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "sub",
		opcode:       opcode17(0b0100000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "sll",
		opcode:       opcode17(0b0000000, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "slt",
		opcode:       opcode17(0b0000000, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "sltu",
		opcode:       opcode17(0b0000000, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "xor",
		opcode:       opcode17(0b0000000, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "srl",
		opcode:       opcode17(0b0000000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "sra",
		opcode:       opcode17(0b0100000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "or",
		opcode:       opcode17(0b0000000, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "and",
		opcode:       opcode17(0b0000000, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    instruction.TypeAritm,
	}, {
		name: "fence",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b0001111}),
			Mask:  revertBytes([]byte{0xf0, 0x0f, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    instruction.TypeMemOrder,
	}, {
		name: "fence.i",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 1 << 4, 0b0001111}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    instruction.TypeMemOrder,
	}, {
		name: "ecall",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    instruction.TypeSyscall,
	}, {
		name: "ebreak",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 1 << 4, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    instruction.TypeSyscall,
	},
	{
		name:                "csrrw",
		opcode:              opcode10(0b001, 0b1110011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		immediate:           immTypeI,
		additionalImmediate: addImmCSR,
		instrType:           instruction.TypeCPUStateChange,
	}, {
		name:                "csrrs",
		opcode:              opcode10(0b010, 0b1110011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		immediate:           immTypeI,
		additionalImmediate: addImmCSR,
		instrType:           instruction.TypeCPUStateChange,
	}, {
		name:                "csrrc",
		opcode:              opcode10(0b011, 0b1110011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		immediate:           immTypeI,
		additionalImmediate: addImmCSR,
		instrType:           instruction.TypeCPUStateChange,
	},
	// FIXME: There are 2 independent immediate values in those 3
	// instructions. Find a way how to parse and represent those
	// instructions.
	{
		name:                "csrrwi",
		opcode:              opcode10(0b101, 0b1110011),
		inputRegCnt:         0,
		hasOutputReg:        true,
		additionalImmediate: addImmCSR,
		instrType:           instruction.TypeCPUStateChange,
	}, {
		name:                "csrrsi",
		opcode:              opcode10(0b110, 0b1110011),
		inputRegCnt:         0,
		hasOutputReg:        true,
		additionalImmediate: addImmCSR,
		instrType:           instruction.TypeCPUStateChange,
	}, {
		name:                "csrrci",
		opcode:              opcode10(0b111, 0b1110011),
		inputRegCnt:         0,
		hasOutputReg:        true,
		additionalImmediate: addImmCSR,
		instrType:           instruction.TypeCPUStateChange,
	},
}

var arithm64 = []*instructionOpcode{
	{
		name:         "lwu",
		opcode:       opcode10(0b110, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    instruction.TypeLoad,
	}, {
		name:         "ld",
		opcode:       opcode10(0b011, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    8,
		immediate:    immTypeI,
		instrType:    instruction.TypeLoad,
	}, {
		name:         "sd",
		opcode:       opcode10(0b011, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   8,
		immediate:    immTypeS,
		instrType:    instruction.TypeStore,
	}, {
		name:                "slli",
		opcode:              opcodeShiftImm(false, 6, 0b001, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           instruction.TypeAritm,
	}, {
		name:                "srli",
		opcode:              opcodeShiftImm(false, 6, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           instruction.TypeAritm,
	}, {
		name:                "srai",
		opcode:              opcodeShiftImm(true, 6, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           instruction.TypeAritm,
	}, {
		name:         "addiw",
		opcode:       opcode10(0b000, 0b0011011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    instruction.TypeAritm,
	}, {
		name:                "slliw",
		opcode:              opcodeShiftImm(false, 5, 0b001, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           instruction.TypeAritm,
	}, {
		name:                "srliw",
		opcode:              opcodeShiftImm(false, 5, 0b101, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           instruction.TypeAritm,
	}, {
		name:                "sraiw",
		opcode:              opcodeShiftImm(true, 5, 0b101, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           instruction.TypeAritm,
	}, {
		name:         "addw",
		opcode:       opcode17(0b0000000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "subw",
		opcode:       opcode17(0b0100000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "sllw",
		opcode:       opcode17(0b0000000, 0b001, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "srlw",
		opcode:       opcode17(0b0000000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "sraw",
		opcode:       opcode17(0b0100000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	},
}

var mul32 = []*instructionOpcode{
	{
		name:         "mul",
		opcode:       opcode17(0b0000001, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "mulh",
		opcode:       opcode17(0b0000001, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "mulhsu",
		opcode:       opcode17(0b0000001, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "mulhu",
		opcode:       opcode17(0b0000001, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "div",
		opcode:       opcode17(0b0000001, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "divu",
		opcode:       opcode17(0b0000001, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "rem",
		opcode:       opcode17(0b0000001, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "remu",
		opcode:       opcode17(0b0000001, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	},
}

var mul64 = []*instructionOpcode{
	{
		name:         "mulw",
		opcode:       opcode17(0b0000001, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "divw",
		opcode:       opcode17(0b0000001, 0b100, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "divuw",
		opcode:       opcode17(0b0000001, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "remw",
		opcode:       opcode17(0b0000001, 0b110, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	}, {
		name:         "remuw",
		opcode:       opcode17(0b0000001, 0b111, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    instruction.TypeAritm,
	},
}

// mergeInstructions merges multiple lists of instruntionOpcode to a single
// list.
func mergeInstructions(lists ...[]*instructionOpcode) []*instructionOpcode {
	length := 0
	for _, a := range lists {
		length += len(a)
	}

	merged := make([]*instructionOpcode, 0, length)
	for _, a := range lists {
		merged = append(merged, a...)
	}

	return merged
}

// overrideInstructions applies X-bit architecture instruction changes to a
// previous version (typically X/2-bit) of instructions based in instruction
// names.
//
// Every X-bit architecture extension has to define some new instructions, but
// more importantly it has to redefine some previous instructions to fit well to
// the new architecture with (typically) twice as big registers. For this
// reason, we need a way how to filter out instructions from the previous
// architecture which has been redefined. We do so simply bby filtering those
// instructions which names match to a name of some instruction in an override
// list.
//
// Filtering based on instruction name is not ideal and works mostly because of
// RISC V convention where an instruction without prefix (i.e. add, sub, sll) is
// used to operate on XLEN bits and there are defined new instructions to
// operate XLEN/2 bit portions of XLEN bit long registers. We could also use
// opcode-wise filtering. But as definition of opcode equivalence is not trivial
// to even define and has some corner cases (for example slli, srli, srai), we
// have decided to avoid opcode code comparison. The logic would be more error
// prone then this simple comparison based in instruction names.
func overrideInstructions(
	instrs []*instructionOpcode,
	overrides []*instructionOpcode,
) []*instructionOpcode {
	// Create own sorted copy to allow quick O(log(n)) search. This trick
	// decreases the overall complexity from O(n^2) to O(n*log(n))
	os := make([]*instructionOpcode, len(overrides))
	copy(os, overrides)
	sort.Slice(os, func(i, j int) bool {
		return strings.Compare(os[i].name, os[j].name) == -1
	})

	// Preallocate the worst possible case - as there will be few filtered
	// instructions, we are not waisting as much as in case of exponential
	// growth of the buffer, which would most likely result in significantly
	// bigger array.
	replaced := make([]*instructionOpcode, 0, len(instrs)+len(overrides))
	for _, instr := range instrs {
		i := sort.Search(len(os), func(i int) bool {
			return strings.Compare(os[i].name, instr.name) > -1
		})

		if i == len(os) || os[i].name != instr.name {
			replaced = append(replaced, instr)
		}
	}

	replaced = append(replaced, overrides...)
	return replaced
}

var (
	known32 = mergeInstructions(arithm32, mul32)
	known64 = overrideInstructions(
		known32,
		mergeInstructions(
			arithm64,
			mul64,
		),
	)
)
