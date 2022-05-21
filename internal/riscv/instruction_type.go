package riscv

import (
	"fmt"
	"mltwist/internal/opcode"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
)

// instructionLen is length of RISC V opcode in bytes.
const instructionLen = 4

// instructionType describes a single RISC-V instruction opcode.
type instructionType struct {
	// name is a symbolic name of an instruction in assembler code.
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

	// immediate describes an immediate value encoding format in an
	// instruction.
	immediate immType

	// instrType is set of instruction types of an opcode.
	instrType model.Type

	// effects is a function which based on specific instruction i evaluates
	// all effects of the given instruction.
	//
	// The array returned is allowed to contain nil expr.Effect values.
	// Those nils will be interpreted as no effects. The reasoning behind
	// nil effects is simplification of handling of writes to x0 register
	// which are effectively defined as having no side effect.
	effects func(i instruction) []expr.Effect
}

// Opcode returns the opcode definition nof a given instruction type.
func (o instructionType) Opcode() opcode.Opcode { return o.opcode }
func (o instructionType) Name() string          { return o.name }

// isPow2 indicates if number is power of 2. This function return true also for
// n=0.
func isPow2(n uint8) bool { return (n & (n - 1)) == 0 }

// validate checks that instructionType description is valid (follows all the
// assumptions the code and the architecture imposes on the struct).
func (o instructionType) validate(xlenBytes uint8) error {
	if o.name == "" {
		return fmt.Errorf("instruction name cannot be empty")
	}
	if err := o.opcode.Validate(); err != nil {
		return fmt.Errorf("invalid opcode description: %w", err)
	}

	if o.inputRegCnt > 2 {
		return fmt.Errorf("too many input registers: %d", o.inputRegCnt)
	}

	if l := o.loadBytes; l > xlenBytes {
		return fmt.Errorf("load is too wide: %d > XLEN(%d)", l, xlenBytes)
	} else if !isPow2(l) {
		return fmt.Errorf("load width is not power of 2: %d", l)
	}

	if s := o.storeBytes; s > xlenBytes {
		return fmt.Errorf("store is too wide: %d > XLEN(%d)", s, xlenBytes)
	} else if !isPow2(s) {
		return fmt.Errorf("store width is not power of 2: %d", s)
	}

	if o.loadBytes > 0 && o.storeBytes > 0 && !o.instrType.MemOrder() {
		return fmt.Errorf("non-atomic instruction can be either load or store")
	}

	if cnt := o.inputRegCnt; !o.instrType.MemOrder() && o.loadBytes > 0 && cnt != 1 {
		return fmt.Errorf("non-atomic load must have 1 input register: %d", cnt)
	} else if !o.instrType.MemOrder() && o.storeBytes > 0 && cnt != 2 {
		return fmt.Errorf("non-atomic store must have 2 input registers: %d", cnt)
	}

	if o.effects == nil {
		return fmt.Errorf("effects function must be always set")
	}

	return nil
}

// validEffects filters nil effects from a list of effects returned by effects
// function.
//
// As RISC-V has register x0 which means to drop any result of the operation,
// it's possible that effects function of an instructionType returns nil
// expr.Effect. On the other hand the interface defined in model.Instruction
// requires all effects to be non-nil. Consequently, we have to filter nil
// expression away from the array.
func (o instructionType) validEffects(i instruction) []expr.Effect {
	effs := o.effects(i)
	if len(effs) == 0 {
		return nil
	}

	effects := make([]expr.Effect, 0, len(effs))
	for _, e := range effs {
		if e != nil {
			effects = append(effects, e)
		}
	}

	return effects
}

// mergeInstructions merges multiple lists of instruntionOpcode into a single
// list.
func mergeInstructions(lists [][]*instructionType) []*instructionType {
	length := 0
	for _, a := range lists {
		length += len(a)
	}

	merged := make([]*instructionType, 0, length)
	for _, a := range lists {
		merged = append(merged, a...)
	}

	return merged
}
