package riscv

import (
	"fmt"
	"mltwist/internal/opcode"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// Variant represents a RISC-V .
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

type Parser struct {
	decoder *opcode.Decoder[*instructionOpcode]
}

// NewParser creates a new RISC-V instruction parser parsing RISC-V architecture
// v with basic integer instruction set and set of extensions specified by exts.
func NewParser(v Variant, exts ...Extension) Parser {
	arch := instructions[v]

	extensions := make([][]*instructionOpcode, 1, len(exts)+1)
	extensions[0] = arch[extI]
	for _, e := range exts {
		if e == extI || e >= extEnd {
			panic(fmt.Sprintf("invlaid extension: %d", e))
		}
		extensions = append(extensions, arch[e])
	}

	instrs := mergeInstructions(extensions)

	decoder, err := opcode.NewDecoder(instrs...)
	if err != nil {
		panic(fmt.Sprintf("unexpected: %s", err.Error()))
	}

	return Parser{
		decoder: decoder,
	}
}

func (p Parser) Parse(addr model.Addr, bytes []byte) (model.Instruction, error) {
	if l := len(bytes); l < instructionLen {
		return model.Instruction{}, fmt.Errorf(
			"bytes are too short (%d) to represent an instruction opcode", l)
	}

	opcode, ok := p.decoder.Match(bytes)
	if !ok {
		return model.Instruction{}, fmt.Errorf(
			"unknown instruction opcode: 0x%x", bytes[:instructionLen])
	}

	instr := newInstruction(addr, bytes, opcode)
	opcodeEffects := opcode.effects(instr)

	var effects []expr.Effect
	if l := len(opcodeEffects); l > 0 {
		effects = make([]expr.Effect, 0, len(opcodeEffects))
		for _, e := range opcodeEffects {
			if e != nil {
				effects = append(effects, e)
			}
		}
	}

	return model.Instruction{
		Type:    opcode.instrType,
		ByteLen: instructionLen,

		Effects: effects,

		Details: instr,
	}, nil
}
