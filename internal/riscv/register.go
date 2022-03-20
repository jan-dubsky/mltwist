package riscv

import "fmt"

const regBits uint8 = 5

// reg represents a register number position in an instruction opcode. Valid
// values are rd (output register), rs1 (input register 1) and rs2 (input
// register 2).
type reg uint8

const (
	// rd is position of output register in an instruction opcode.
	rd reg = iota + 1
	// rs1 is position of first input register in an instruction opcode.
	rs1
	// rs2 is position of second input register in an instruction opcode.
	rs2
)

// bitOffset returns index if starting bit of register position in in
// instruction opcode. Lowest bit has index 0.
func (r reg) bitOffset() uint8 {
	switch r {
	case rd:
		return 7
	case rs1:
		return 15
	case rs2:
		return 20
	default:
		panic(fmt.Sprintf("invalid register: %d", r))
	}
}

// regNum parses a register number at a given position from b.
func (r reg) regNum(value uint32) regNum {
	const regNumMask = uint32(1)<<regBits - 1

	shifted := value >> r.bitOffset()
	masked := shifted & regNumMask

	return regNum(masked)
}
