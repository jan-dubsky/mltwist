package parser

import "decomp/internal/instruction"

type Strategy interface {
	Window() uint64

	Parse(bytes []byte) (instruction.Instruction, error)
}

// Memory represents memory of a userspace process as it would look at the time
// of userspace program start.
//
// The memory is allowed to be sparse, as even the process memory can be sparse.
// On top of that, it's allowed for the memory not to represent the exact memory
// space of a userspace process, but it's allowed to omit some parts of the
// memory. This is typically useful when for example ELF file contains also
// symbols for debugging, exceptions etc. But it's also allowed to omit some
// relevant parts of the memory as for example static data.
type Memory interface {
	// Bytes returns a sequence of length bytes from program memory starting
	// at position offset.
	//
	// The slice returned can refer the internal representation of memory,
	// so modification of the slice concent might change the internal state.
	// For this reason, the returned slice has to be treated as read-only.
	//
	// As memory is allowed to be sparse, it's possible that some bytes in
	// range [offset, offset + length) won't exist in the sparse memory. In
	// such a case this method has to indicate such situation by returning
	// nil slice.
	Bytes(offset uint64, length uint64) []byte
}
