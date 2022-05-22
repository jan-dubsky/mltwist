package riscv

import (
	"fmt"
	"mltwist/internal/opcode"
	"mltwist/pkg/model"
)

// Variant represents a variant of RISC-V address and register width.
type Variant uint8

const (
	// Variant32 represents RISC-V 32 bit architecture - rv32.
	Variant32 Variant = iota
	// Variant64 represents RISC-V 64 bit architecture - rv64.
	Variant64

	// extEnd marks first invalid value of architecture.
	variantEnd
)

// Modules encodes RISC-V extension/module according to RISC-V ISA (I, M, A,
// etc.).
type Extension uint8

const (
	// ExtI represents integer operations basic instruction set of RISC-V
	// ISA.
	//
	// This extension is automatically included in every parser as I base
	// instruction set is the basic extension introducing necessary
	// instruction in the RISC-V instruction set.
	extI Extension = iota

	// ExtM represents integer multiplication and division extension of
	// RISC-V ISA.
	ExtM

	// ExtA represents atomic instruction extension of RISC-V ISA.
	ExtA

	// extEnd marks first invalid value of extension.
	extEnd
)

// Parser parses RISC-V instructions of a specified variant with specified set
// of extensions.
type Parser struct {
	matcher *opcode.Matcher[*instructionType]
}

// NewParser creates a new RISC-V instruction parser parsing RISC-V architecture
// v with basic integer instruction set and set of extensions specified by exts.
//
// Usage variant not exported by this package is undefined. The same applies to
// extensions not specified by this package. Specifying one extension multiple
// times is undefined as well.
func NewParser(v Variant, exts ...Extension) Parser {
	instrs := instructionSet(v, exts)
	decoder, err := opcode.NewMatcher(instrs)

	// This means that instruction opcodes defined in this package are
	// either invalid or they collide. This is non-recoverable as the
	// package code has to be modified.
	if err != nil {
		panic(fmt.Sprintf("bug: matcher creation failed: %s", err.Error()))
	}

	return Parser{
		matcher: decoder,
	}
}

// instructionSet generates a new set of instructions based on RISC-V variant v
// and list of extensions. The extI extension instructions are always added to
// the set.
func instructionSet(v Variant, exts []Extension) []*instructionType {
	variantExtensions := instructions[v]

	extensions := make([][]*instructionType, 1, len(exts)+1)
	extensions[0] = variantExtensions[extI]
	for _, e := range exts {
		// This is technically a user error which could be recoverable.
		// But it's violation of API contract, so panic is in a way
		// appropriate punishment.
		if e == extI || e >= extEnd {
			panic(fmt.Sprintf("invalid extension: %d", e))
		}

		// The map is written in the code -> changing this requires code
		// modification. As this issue is non-recoverable, panic it is.
		ext, ok := variantExtensions[e]
		if !ok {
			panic(fmt.Sprintf("bug: valid extension not found: %d", e))
		}

		extensions = append(extensions, ext)
	}

	return mergeInstructions(extensions)
}

// Parse parses an instruction starting at address a comprising of bytes at the
// beginning of bs.
//
// The array of bytes bs is allowed to be longer than the instruction. In such a
// case the Parse method will take into consideration only those bytes at the
// beginning of the array which represent a single instruction. This allows to
// implement paginated parsing of a sequence of instructions by incrementally
// cutting the start of bs.
func (p Parser) Parse(a model.Addr, bs []byte) (model.Instruction, error) {
	if l := len(bs); l < instructionLen {
		return model.Instruction{}, fmt.Errorf(
			"bytes are too short to be a RISCV instruction opcode: %d", l)
	}

	opcode, ok := p.matcher.Match(bs)
	if !ok {
		return model.Instruction{}, fmt.Errorf(
			"unknown instruction opcode: 0x%x", bs[:instructionLen])
	}

	instr := newInstruction(a, bs, opcode)
	return model.Instruction{
		Type:    opcode.instrType,
		ByteLen: instructionLen,

		Effects: opcode.validEffects(instr),
		Details: instr,
	}, nil
}
