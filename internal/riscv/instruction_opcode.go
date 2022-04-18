package riscv

import (
	"fmt"
	"mltwist/internal/opcode"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// instructionLen is length of RISC V opcode in bytes.
const instructionLen = 4

// instructionOpcode describes a single RISC-V instruction opcode.
type instructionOpcode struct {
	// name is a symbolic name of an instruction in assember code.
	//
	// For example: lui, bl, lb, sw etc.
	name string
	// opcode describes opcode bits in an instruction.
	opcode opcode.Opcode

	// inputRegCnt is number of input registers of an instruction. Valid
	// values are 0, 1 and 2.
	inputRegCnt uint8
	// hasOutputReg indicates whether an instruction writes output to
	// register.
	hasOutputReg bool

	// loadBytes is number of bytes the instruction loads from memory.
	loadBytes uint8
	// storeBytes is number of bytes the instruction stores into memory.
	storeBytes uint8

	// unsigned marks if load from memory which is shorter than XLEN bits is
	// unsigned or signed (default).
	unsigned bool

	// immediate describes an immediate value encoding format in an
	// instruction.
	immediate           immType
	additionalImmediate addOpcodeInfo

	// instrType is set of instruction types of an opcode.
	instrType model.Type

	effects func(i Instruction) []expr.Effect
}

func (i instructionOpcode) Opcode() opcode.Opcode { return i.opcode }
func (i instructionOpcode) Name() string          { return i.name }

// isPow2 indicates if number is power of 2. This function return true also for
// n=0.
func isPow2(n uint8) bool { return (n & (n - 1)) == 0 }

// validate checks that instructionOpcode is valid (follows all the assumptions
// the code and the architecture imposes on the struct).
func (i instructionOpcode) validate(xlenBytes uint8) error {
	if i.name == "" {
		return fmt.Errorf("instruction name cannot be empty")
	}
	if err := i.opcode.Validate(); err != nil {
		return fmt.Errorf("invalid opcode description: %w", err)
	}

	if i.inputRegCnt > 2 {
		return fmt.Errorf("too many input registers: %d", i.inputRegCnt)
	}

	if l := i.loadBytes; l > xlenBytes {
		return fmt.Errorf("load is too wide: %d > XLEN(%d)", l, xlenBytes)
	} else if !isPow2(l) {
		return fmt.Errorf("load width is not power of 2: %d", l)
	}

	if s := i.storeBytes; s > xlenBytes {
		return fmt.Errorf("store is too wide: %d > XLEN(%d)", s, xlenBytes)
	} else if !isPow2(s) {
		return fmt.Errorf("store width is not power of 2: %d", s)
	}

	if i.loadBytes > 0 && i.storeBytes > 0 {
		return fmt.Errorf("instruction can be either load or store")
	}
	if i.unsigned && (i.loadBytes == 0 || i.loadBytes == xlenBytes) {
		return fmt.Errorf("unsigned is allowed for loads shorter than XLEN bytes")
	}

	if cnt := i.inputRegCnt; i.loadBytes > 0 && cnt != 1 {
		return fmt.Errorf("load must have exactly one input register: %d", cnt)
	} else if i.storeBytes > 0 && cnt != 2 {
		return fmt.Errorf("store must have exactly two input registers: %d", cnt)
	}

	if i.effects == nil {
		return fmt.Errorf("effects function must be always set")
	}

	return nil
}

// mergeInstructions merges multiple lists of instruntionOpcode into a single
// list.
func mergeInstructions(lists [][]*instructionOpcode) []*instructionOpcode {
	length := 0
	for _, a := range lists {
		length += len(a)
	}

	merged := make([]*instructionOpcode, 0, length)
	for _, a := range lists {
		merged = append(merged, a...)
	}

	return merged
}
