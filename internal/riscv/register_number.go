package riscv

import "fmt"

// regCnt is number or registers in RISC-V platform.
const regCnt = 32

// regNum represents a RISC-V register number. Range of valid values is [0..31]
// (i.e. [0..regCnt-1]).
type regNum uint8

func (r regNum) String() string { return fmt.Sprintf("r%d", r) }
