package riscv

import "fmt"

const regBits uint8 = 5

type reg uint8

const (
	rd reg = iota + 1
	rs1
	rs2
)

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
