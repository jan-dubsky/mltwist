package riscv

import (
	"decomp/internal/opcode"
	"decomp/pkg/model"
	"fmt"
)

const (
	// low7Bits is a byte (mask) with bottom 7 bits set and last bit unset.
	low7Bits byte = 0x7F
	// low3Bits is a byte (mask) with bottom 3 bits set and all higher bits
	// unset.
	low3Bits byte = 0x7
)

// assertMask checks that only bits set in mask are set in b. This method will
// panic if any other bit is set on b.
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

// opcode7 returns opcode matching low with bottom 7 bits of an instruction.
//
// As 1 byte opcodes have only 7 bits, this method will panic for values of low
// greater than 127.
func opcode7(low byte) opcode.Opcode {
	assertMask(low, low7Bits)

	return opcode.Opcode{
		Bytes: []byte{low},
		Mask:  []byte{low7Bits},
	}
}

// opcode10 returns opcode matching low with bottom 7 bits and mid with bits
// [12..14] of an instruction.
//
// This method panics if low is greater than 127 or if mid is greater than 7.
func opcode10(mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)

	return opcode.Opcode{
		Bytes: []byte{low, mid << 4},
		Mask:  []byte{low7Bits, low3Bits << 4},
	}
}

// opcode10 returns opcode matching low with bottom 7 bits, mid with bits
// [12..14] and high with bits [25..31] of an instruction.
//
// This method panics if low is greater than 127, mid is greater than 7 or if
// high is greater than 127.
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
// Even though RISC V manual states that there are only 6 distinct instruction
// encodings and all of them should be describable by either opcode7, opcode10,
// or opcode17, there is one small exception. Yes x86, you are not the only one
// architecture doing weird things... The exceptional opcode encoding is the one
// with fixed bit short argument in the instruction immediate value.
// Technically, such an instruction is just a bit shift with 12 bit immediate,
// but there are a few catches.
//
// The first catch is that not all immediate values are allowed. To be more
// specific, on an architecture with XLEN bits in registers (for simplicity
// let's consider XLEN=32 - 32 bit processors), it doesn't make sense to encode
// more than 31 bit immediate value to shift and though all (but one - see
// below) higher bits of immediate value are reserved to be zero.
//
// Another irregularity in immediate shift instruction encoding is the different
// in between logical and arithmetic shift. The bit differentiating logical and
// arithmetic shift is bit [30] of an instruction opcode which would correcpond
// to bit [11] of 12bit I immediate type encoding. Unfortunately this encodings
// brings a weird inconsistency when two distinct instructions identified by two
// different assembler names (srli and srai) have the same opcode but differ
// only in an immediate value bit.
//
// As different immediate value can encode different assembler instructions, we
// need them to be parsed as 2 different instructions. Consequently we are
// forced to describe this instruction opcode meta-format which is not specified
// by the architecture specification document by itself, but which allows us to
// parse bit shifts by an immediate value.
//
// The problem with differentiating logical and arithmetic shirt applies as well
// on srl and sra instructions (i.e. shift instructions accepting register
// arguments). Fortunately there we can treat the instruction opcode as 17bit
// opcode as every other bit (but bit [30]) of an immediate value is reserved to
// be zero.
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
	// to ensure that all other reserved bits of the actual opcode are zero.
	shiftBitMask := (uint16(1) << shiftBits) - 1
	// Then we have to shift this mask to the right place - to 20th bit of
	// opcode. As we have just high half of instruction opcode, we are
	// already shifted 16 bits. So 4 bits are remaining.
	highHalfMask := (^shiftBitMask) << 4

	return opcode.Opcode{
		Bytes: []byte{low, mid << 4, 0, high},
		Mask: []byte{
			low7Bits,
			low3Bits << 4,
			byte(highHalfMask),
			byte(highHalfMask >> 8),
		},
	}
}

var integer32 = []*instructionOpcode{
	{
		name:         "lui",
		opcode:       opcode7(0b0110111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    model.TypeAritm,
	}, {
		name:         "auipc",
		opcode:       opcode7(0b0010111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    model.TypeAritm,
	}, {
		name:         "jal",
		opcode:       opcode7(0b1101111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeJ,
		instrType:    model.TypeJump,
	}, {
		name:         "jalr",
		opcode:       opcode10(0b000, 0b1100111),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeJumpDyn,
	}, {
		name:         "beq",
		opcode:       opcode10(0b000, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
	}, {
		name:         "bne",
		opcode:       opcode10(0b001, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
	}, {
		name:         "blt",
		opcode:       opcode10(0b100, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
	}, {
		name:         "bge",
		opcode:       opcode10(0b101, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
	}, {
		name:         "bltu",
		opcode:       opcode10(0b110, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
	}, {
		name:         "bgeu",
		opcode:       opcode10(0b111, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
	}, {
		name:         "lb",
		opcode:       opcode10(0b000, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
	}, {
		name:         "lh",
		opcode:       opcode10(0b001, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
	}, {
		name:         "lw",
		opcode:       opcode10(0b010, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
	}, {
		name:         "lbu",
		opcode:       opcode10(0b100, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
	}, {
		name:         "lhu",
		opcode:       opcode10(0b101, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
	}, {
		name:         "sb",
		opcode:       opcode10(0b000, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   1,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
	}, {
		name:         "sh",
		opcode:       opcode10(0b001, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   2,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
	}, {
		name:         "sw",
		opcode:       opcode10(0b010, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   4,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
	}, {
		name:         "addi",
		opcode:       opcode10(0b000, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
	}, {
		name:         "slti",
		opcode:       opcode10(0b010, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
	}, {
		name:         "sltiu",
		opcode:       opcode10(0b011, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
	}, {
		name:         "xori",
		opcode:       opcode10(0b100, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
	}, {
		name:         "ori",
		opcode:       opcode10(0b110, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
	}, {
		name:         "andi",
		opcode:       opcode10(0b111, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
	}, {
		name:                "slli",
		opcode:              opcodeShiftImm(false, 5, 0b001, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
	}, {
		name:                "srli",
		opcode:              opcodeShiftImm(false, 5, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
	}, {
		name:                "srai",
		opcode:              opcodeShiftImm(true, 5, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
	}, {
		name:         "add",
		opcode:       opcode17(0b0000000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "sub",
		opcode:       opcode17(0b0100000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "sll",
		opcode:       opcode17(0b0000000, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "slt",
		opcode:       opcode17(0b0000000, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "sltu",
		opcode:       opcode17(0b0000000, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "xor",
		opcode:       opcode17(0b0000000, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "srl",
		opcode:       opcode17(0b0000000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "sra",
		opcode:       opcode17(0b0100000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "or",
		opcode:       opcode17(0b0000000, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name:         "and",
		opcode:       opcode17(0b0000000, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
	}, {
		name: "fence",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b0001111}),
			Mask:  revertBytes([]byte{0xf0, 0x0f, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeMemOrder,
	}, {
		name: "fence.i",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 1 << 4, 0b0001111}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeMemOrder,
	}, {
		name: "ecall",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
	}, {
		name: "ebreak",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 1 << 4, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
	},
	{
		name:                "csrrw",
		opcode:              opcode10(0b001, 0b1110011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		immediate:           immTypeI,
		additionalImmediate: addImmCSR,
		instrType:           model.TypeCPUStateChange,
	}, {
		name:                "csrrs",
		opcode:              opcode10(0b010, 0b1110011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		immediate:           immTypeI,
		additionalImmediate: addImmCSR,
		instrType:           model.TypeCPUStateChange,
	}, {
		name:                "csrrc",
		opcode:              opcode10(0b011, 0b1110011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		immediate:           immTypeI,
		additionalImmediate: addImmCSR,
		instrType:           model.TypeCPUStateChange,
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
		instrType:           model.TypeCPUStateChange,
	}, {
		name:                "csrrsi",
		opcode:              opcode10(0b110, 0b1110011),
		inputRegCnt:         0,
		hasOutputReg:        true,
		additionalImmediate: addImmCSR,
		instrType:           model.TypeCPUStateChange,
	}, {
		name:                "csrrci",
		opcode:              opcode10(0b111, 0b1110011),
		inputRegCnt:         0,
		hasOutputReg:        true,
		additionalImmediate: addImmCSR,
		instrType:           model.TypeCPUStateChange,
	},
}

var integer64 = []*instructionOpcode{
	{
		name:         "lwu",
		opcode:       opcode10(0b110, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
	}, {
		name:         "ld",
		opcode:       opcode10(0b011, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    8,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
	}, {
		name:         "sd",
		opcode:       opcode10(0b011, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   8,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
	}, {
		name:                "slli",
		opcode:              opcodeShiftImm(false, 6, 0b001, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           model.TypeAritm,
	}, {
		name:                "srli",
		opcode:              opcodeShiftImm(false, 6, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           model.TypeAritm,
	}, {
		name:                "srai",
		opcode:              opcodeShiftImm(true, 6, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           model.TypeAritm,
	}, {
		name:         "addiw",
		opcode:       opcode10(0b000, 0b0011011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
	}, {
		name:                "slliw",
		opcode:              opcodeShiftImm(false, 5, 0b001, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
	}, {
		name:                "srliw",
		opcode:              opcodeShiftImm(false, 5, 0b101, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
	}, {
		name:                "sraiw",
		opcode:              opcodeShiftImm(true, 5, 0b101, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
	}, {
		name:         "addw",
		opcode:       opcode17(0b0000000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "subw",
		opcode:       opcode17(0b0100000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "sllw",
		opcode:       opcode17(0b0000000, 0b001, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "srlw",
		opcode:       opcode17(0b0000000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "sraw",
		opcode:       opcode17(0b0100000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	},
}

var mul32 = []*instructionOpcode{
	{
		name:         "mul",
		opcode:       opcode17(0b0000001, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "mulh",
		opcode:       opcode17(0b0000001, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "mulhsu",
		opcode:       opcode17(0b0000001, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "mulhu",
		opcode:       opcode17(0b0000001, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "div",
		opcode:       opcode17(0b0000001, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "divu",
		opcode:       opcode17(0b0000001, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "rem",
		opcode:       opcode17(0b0000001, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "remu",
		opcode:       opcode17(0b0000001, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	},
}

var mul64 = []*instructionOpcode{
	{
		name:         "mulw",
		opcode:       opcode17(0b0000001, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "divw",
		opcode:       opcode17(0b0000001, 0b100, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "divuw",
		opcode:       opcode17(0b0000001, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "remw",
		opcode:       opcode17(0b0000001, 0b110, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	}, {
		name:         "remuw",
		opcode:       opcode17(0b0000001, 0b111, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
	},
}

var instructions = map[Variant]map[Extension][]*instructionOpcode{
	Variant32: {
		extI: integer32,
		ExtM: mul32,
	},
	Variant64: {
		extI: overrideInstructions(integer32, integer64),
		ExtM: overrideInstructions(mul32, mul64),
	},
}
