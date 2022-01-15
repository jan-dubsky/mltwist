package riscv

import (
	"decomp/internal/opcode"
	"fmt"
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

func opcode7(opc byte) opcode.Opcode {
	assertMask(opc, low7Bits)

	return opcode.Opcode{
		Bytes: []byte{opc},
		Mask:  []byte{low7Bits},
	}
}

func opcode10(high byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(high, low3Bits)

	return opcode.Opcode{
		Bytes: []byte{high << 4, low},
		Mask:  []byte{low3Bits << 4, low7Bits},
	}
}

func opcode17(high byte, mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)
	assertMask(high, low7Bits)

	return opcode.Opcode{
		Bytes: []byte{high << 1, 0, mid << 4, low},
		Mask:  []byte{low7Bits << 1, 0, low3Bits << 4, low7Bits},
	}
}

var arithm32 = []*instr{
	{
		name:   "lui",
		opcode: opcode7(0b0110111),
	}, {
		name:   "auipc",
		opcode: opcode7(0b0010111),
	}, {
		name:   "jal",
		opcode: opcode7(0b1101111),
	}, {
		name:   "jalr",
		opcode: opcode10(0b000, 0b1100111),
	}, {
		name:   "beq",
		opcode: opcode10(0b000, 0b1100011),
	}, {
		name:   "bne",
		opcode: opcode10(0b001, 0b1100011),
	}, {
		name:   "blt",
		opcode: opcode10(0b100, 0b1100011),
	}, {
		name:   "bge",
		opcode: opcode10(0b101, 0b1100011),
	}, {
		name:   "bltu",
		opcode: opcode10(0b110, 0b1100011),
	}, {
		name:   "bgeu",
		opcode: opcode10(0b111, 0b1100011),
	}, {
		name:   "lb",
		opcode: opcode10(0b000, 0b0000011),
	}, {
		name:   "lh",
		opcode: opcode10(0b001, 0b0000011),
	}, {
		name:   "lw",
		opcode: opcode10(0b010, 0b0000011),
	}, {
		name:   "lbu",
		opcode: opcode10(0b100, 0b0000011),
	}, {
		name:   "lhu",
		opcode: opcode10(0b101, 0b0000011),
	}, {
		name:   "sb",
		opcode: opcode10(0b000, 0b0100011),
	}, {
		name:   "sh",
		opcode: opcode10(0b001, 0b0100011),
	}, {
		name:   "sw",
		opcode: opcode10(0b010, 0b0100011),
	}, {
		name:   "addi",
		opcode: opcode10(0b000, 0b0010011),
	}, {
		name:   "slti",
		opcode: opcode10(0b010, 0b0010011),
	}, {
		name:   "sltiu",
		opcode: opcode10(0b011, 0b0010011),
	}, {
		name:   "xori",
		opcode: opcode10(0b100, 0b0010011),
	}, {
		name:   "ori",
		opcode: opcode10(0b110, 0b0010011),
	}, {
		name:   "andi",
		opcode: opcode10(0b111, 0b0010011),
	}, {
		name:   "slli",
		opcode: opcode17(0b0000000, 0b001, 0b0010011),
	}, {
		name:   "srli",
		opcode: opcode17(0b0000000, 0b101, 0b0010011),
	}, {
		name:   "srai",
		opcode: opcode17(0b0100000, 0b101, 0b0010011),
	}, {
		name:   "add",
		opcode: opcode17(0b0000000, 0b000, 0b0110011),
	}, {
		name:   "sub",
		opcode: opcode17(0b0100000, 0b000, 0b0110011),
	}, {
		name:   "sll",
		opcode: opcode17(0b0000000, 0b001, 0b0110011),
	}, {
		name:   "slt",
		opcode: opcode17(0b0000000, 0b010, 0b0110011),
	}, {
		name:   "sltu",
		opcode: opcode17(0b0000000, 0b011, 0b0110011),
	}, {
		name:   "xor",
		opcode: opcode17(0b0000000, 0b100, 0b0110011),
	}, {
		name:   "srl",
		opcode: opcode17(0b0000000, 0b101, 0b0110011),
	}, {
		name:   "sra",
		opcode: opcode17(0b0100000, 0b101, 0b0110011),
	}, {
		name:   "or",
		opcode: opcode17(0b0000000, 0b110, 0b0110011),
	}, {
		name:   "and",
		opcode: opcode17(0b0000000, 0b111, 0b0110011),
	}, {
		name: "fence",
		opcode: opcode.Opcode{
			Bytes: []byte{0, 0, 0, 0b0001111},
			Mask:  []byte{0xf0, 0x0f, 0xff, 0xff},
		},
	}, {
		name: "fence.i",
		opcode: opcode.Opcode{
			Bytes: []byte{0, 0, 1 << 4, 0b0001111},
			Mask:  []byte{0xff, 0xff, 0xff, 0xff},
		},
	}, {
		name: "ecall",
		opcode: opcode.Opcode{
			Bytes: []byte{0, 0, 0, 0b1110011},
			Mask:  []byte{0xff, 0xff, 0xff, 0xff},
		},
	}, {
		name: "ebreak",
		opcode: opcode.Opcode{
			Bytes: []byte{0, 1 << 4, 0, 0b1110011},
			Mask:  []byte{0xff, 0xff, 0xff, 0xff},
		},
	}, {
		name:   "csrrw",
		opcode: opcode10(0b001, 0b1110011),
	}, {
		name:   "csrrs",
		opcode: opcode10(0b010, 0b1110011),
	}, {
		name:   "csrrc",
		opcode: opcode10(0b011, 0b1110011),
	}, {
		name:   "csrrwi",
		opcode: opcode10(0b101, 0b1110011),
	}, {
		name:   "csrrsi",
		opcode: opcode10(0b110, 0b1110011),
	}, {
		name:   "csrrci",
		opcode: opcode10(0b111, 0b1110011),
	},
}

var arithm64 = []*instr{
	{
		name:   "lwu",
		opcode: opcode10(0b110, 0b0000011),
	}, {
		name:   "ld",
		opcode: opcode10(0b011, 0b0000011),
	}, {
		name:   "sd",
		opcode: opcode10(0b011, 0b0100011),
	}, {
		name:   "slli",
		opcode: opcode17(0b0000000, 0b001, 0b0010011),
	}, {
		name:   "srli",
		opcode: opcode17(0b0000000, 0b101, 0b0010011),
	}, {
		name:   "srai",
		opcode: opcode17(0b010000, 0b101, 0b0010011),
	}, {
		name:   "addiw",
		opcode: opcode10(0b000, 0b0011011),
	}, {
		name:   "slliw",
		opcode: opcode17(0b0000000, 0b001, 0b0011011),
	}, {
		name:   "srliw",
		opcode: opcode17(0b0000000, 0b101, 0b0011011),
	}, {
		name:   "sraiw",
		opcode: opcode17(0b0100000, 0b101, 0b0011011),
	}, {
		name:   "addw",
		opcode: opcode17(0b0000000, 0b000, 0b0111011),
	}, {
		name:   "subw",
		opcode: opcode17(0b0100000, 0b000, 0b0111011),
	}, {
		name:   "sllw",
		opcode: opcode17(0b0000000, 0b001, 0b0111011),
	}, {
		name:   "srlw",
		opcode: opcode17(0b0000000, 0b101, 0b0111011),
	}, {
		name:   "sraw",
		opcode: opcode17(0b0100000, 0b101, 0b0111011),
	},
}

var mul32 = []*instr{
	{
		name:   "mul",
		opcode: opcode17(0b0000001, 0b000, 0b0110011),
	}, {
		name:   "mulh",
		opcode: opcode17(0b0000001, 0b001, 0b0110011),
	}, {
		name:   "mulhsu",
		opcode: opcode17(0b0000001, 0b010, 0b0110011),
	}, {
		name:   "mulhu",
		opcode: opcode17(0b0000001, 0b011, 0b0110011),
	}, {
		name:   "div",
		opcode: opcode17(0b0000001, 0b100, 0b0110011),
	}, {
		name:   "divu",
		opcode: opcode17(0b0000001, 0b101, 0b0110011),
	}, {
		name:   "rem",
		opcode: opcode17(0b0000001, 0b110, 0b0110011),
	}, {
		name:   "remu",
		opcode: opcode17(0b0000001, 0b111, 0b0110011),
	},
}

var mul64 = []*instr{
	{
		name:   "mulw",
		opcode: opcode17(0b0000001, 0b000, 0b0111011),
	}, {
		name:   "divw",
		opcode: opcode17(0b0000001, 0b100, 0b0111011),
	}, {
		name:   "divuw",
		opcode: opcode17(0b0000001, 0b101, 0b0111011),
	}, {
		name:   "remw",
		opcode: opcode17(0b0000001, 0b110, 0b0111011),
	}, {
		name:   "remuw",
		opcode: opcode17(0b0000001, 0b111, 0b0111011),
	},
}

func mergeInstructions(arrays ...[]*instr) []*instr {
	length := 0
	for _, a := range arrays {
		length += len(a)
	}

	merged := make([]*instr, 0, length)
	for _, a := range arrays {
		merged = append(merged, a...)
	}

	return merged
}

var (
	known32 = mergeInstructions(arithm32, mul32)
	known64 = mergeInstructions(
		arithm32,
		arithm64,
		mul32,
		mul64,
	)
)
