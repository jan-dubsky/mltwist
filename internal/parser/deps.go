package parser

import "mltwist/pkg/model"

// Parser is a platform specific object responsible for instruction parsing.
type Parser interface {
	// Parse matches an instruction opcode, identifies a platform-specific
	// instruction are returns its generic representation.
	//
	// It's guaranteed that b will start at an instruction boundary and that
	// start of b will be at address addr. On the other hand length of b is
	// not limited. The motivation behind this decision is to support
	// platforms with non-constant instruction length.
	//
	// Byte slice b has to be treated as read-only. Any write or
	// modification of b is considered an API violation.
	//
	// This method can fail in case it's not possible to identify an
	// instruction, b doesn't start with a valid opcode or b is too short to
	// contain the full instruction opcode.
	Parse(addr model.Addr, b []byte) (model.Instruction, error)
}
