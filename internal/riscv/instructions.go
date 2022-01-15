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

var known = []*instr{
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
