package riscv

import "fmt"

// csr represents an arbitrary control and status register. Only bottom 12 bits
// of this value are valid.
type csr uint16

func (c csr) String() string { return fmt.Sprintf("csr%d", c) }
