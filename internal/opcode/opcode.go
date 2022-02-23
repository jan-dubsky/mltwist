package opcode

import "fmt"

// OpcodeGetter is an arbitrary type which exposes an opcode.
type OpcodeGetter interface {
	// Opcode returns the opcode definition.
	Opcode() Opcode
	// Name returns a name of an instruction the Opcode belong to.
	Name() string
}

// Opcode represents a sequence of bits representing an opcode in any possible
// ISA.
type Opcode struct {
	// Bytes are bytes representing an instruction opcode.
	//
	// See Mask to see how to encode non-continuous opcode bytes or
	// individual opcode bits.
	Bytes []byte

	// Mask is a bitmask of bits of Bytes which are considered valid.
	//
	// As there might be ISAs which do not use full bytes or non-continuous
	// sequence of bytes to encode instruction opcode, we need a way how to
	// describe which bits of Bytes are part of opcode and which bits have
	// different meanint then opcode.
	//
	// This array must have the same length as Bytes.
	//
	// Out of Mask nature, it makes no sense for its last byte to be zero as
	// then, both Bytes and Mask could be one byte shorter. For this reason
	// a valid mask must end with nonzero byte.
	Mask []byte
}

func (o Opcode) Validate() error {
	if len(o.Bytes) == 0 {
		return fmt.Errorf("opcode has zero bytes")
	}

	if len(o.Bytes) != len(o.Mask) {
		return fmt.Errorf("bytes length (%d) differs from mask length (%d)",
			len(o.Bytes), len(o.Mask))
	}

	if o.Mask[len(o.Mask)-1] == 0 {
		return fmt.Errorf("last byte of the mask is zero: %v", o.Mask)
	}

	return nil
}

// String returns a human readable representation of an opcode.
func (o Opcode) String() string {
	return fmt.Sprintf("bytes: 0x%x (mask: 0x%x)", o.Bytes, o.Mask)
}
