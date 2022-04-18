package riscv

import (
	"mltwist/internal/opcode"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"mltwist/pkg/model"
)

var integer64 = []*instructionOpcode{
	{
		name:         "lui",
		opcode:       opcode7(0b0110111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			imm, _ := immTypeU.parseValue(i.value)
			val := sext32To64(expr.NewConstInt(imm, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "auipc",
		opcode:       opcode7(0b0010111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			imm, _ := immTypeU.parseValue(i.value)
			val := expr.NewConstUint(addrAddImm(i.address, imm), width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "jal",
		opcode:       opcode7(0b1101111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeJ,
		instrType:    model.TypeJump,
		effects: func(i Instruction) []expr.Effect {
			target := addrImmConst(immTypeJ, i, width64)
			// Address of following instruction.
			following := expr.NewConstUint(i.address+4, width64)
			return []expr.Effect{
				expr.NewRegStore(target, expr.IPKey, width64),
				regStore(following, i, width64),
			}
		},
	}, {
		name:         "jalr",
		opcode:       opcode10(0b000, 0b1100111),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeJumpDyn,
		// FIXME: Find a way how to represent those jump targets.
		effects: func(i Instruction) []expr.Effect {
			target := regImmOp(expr.Add, immTypeI, i, width64)
			// Address of following instruction.
			following := expr.NewConstUint(i.address+4, width64)
			return []expr.Effect{
				expr.NewRegStore(target, expr.IPKey, width64),
				regStore(following, i, width64),
			}
		},
	}, {
		name:         "beq",
		opcode:       opcode10(0b000, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Eq, true, i, width64)}
		},
	}, {
		name:         "bne",
		opcode:       opcode10(0b001, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Eq, false, i, width64)}
		},
	}, {
		name:         "blt",
		opcode:       opcode10(0b100, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Lts, true, i, width64)}
		},
	}, {
		name:         "bge",
		opcode:       opcode10(0b101, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Lts, false, i, width64)}
		},
	}, {
		name:         "bltu",
		opcode:       opcode10(0b110, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Ltu, true, i, width64)}
		},
	}, {
		name:         "bgeu",
		opcode:       opcode10(0b111, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Ltu, false, i, width64)}
		},
	}, {
		name:         "lb",
		opcode:       opcode10(0b000, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width64)
			val := sext(memLoad(addr, width8), 7, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "lh",
		opcode:       opcode10(0b001, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width64)
			val := sext(memLoad(addr, width16), 15, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "lw",
		opcode:       opcode10(0b010, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width64)
			val := sext(memLoad(addr, width32), 31, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "ld",
		opcode:       opcode10(0b011, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    8,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width64)
			val := memLoad(addr, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "lbu",
		opcode:       opcode10(0b100, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width64)
			val := exprtools.Crop(memLoad(addr, width8), width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "lhu",
		opcode:       opcode10(0b101, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width64)
			val := exprtools.Crop(memLoad(addr, width16), width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "lwu",
		opcode:       opcode10(0b110, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		unsigned:     true,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width64)
			val := exprtools.Crop(memLoad(addr, width32), width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sb",
		opcode:       opcode10(0b000, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   1,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(expr.Add, immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width8)}
		},
	}, {
		name:         "sh",
		opcode:       opcode10(0b001, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   2,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(expr.Add, immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width16)}
		},
	}, {
		name:         "sw",
		opcode:       opcode10(0b010, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   4,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(expr.Add, immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width32)}
		},
	}, {
		name:         "sd",
		opcode:       opcode10(0b011, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   8,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(expr.Add, immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width64)}
		},
	}, {
		name:         "addi",
		opcode:       opcode10(0b000, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.Add, immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "slti",
		opcode:       opcode10(0b010, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Lts,
				regLoad(rs1, i, width64),
				immConst(immTypeI, i),
				expr.One,
				expr.Zero,
				width64,
			)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sltiu",
		opcode:       opcode10(0b011, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Ltu,
				regLoad(rs1, i, width64),
				immConst(immTypeI, i),
				expr.One,
				expr.Zero,
				width64,
			)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "xori",
		opcode:       opcode10(0b100, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.Xor, immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "ori",
		opcode:       opcode10(0b110, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.Or, immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "andi",
		opcode:       opcode10(0b111, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.And, immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:                "slli",
		opcode:              opcodeShiftImm(false, 6, 0b001, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := regImmShift(expr.Lsh, i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:                "srli",
		opcode:              opcodeShiftImm(false, 6, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := regImmShift(expr.Rsh, i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:                "srai",
		opcode:              opcodeShiftImm(true, 6, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh64,
		instrType:           model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			reg := regLoad(rs1, i, width64)
			val := exprtools.RshA(reg, immShift(6, i), width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "add",
		opcode:       opcode17(0b0000000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Add, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sub",
		opcode:       opcode17(0b0100000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Sub, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "slt",
		opcode:       opcode17(0b0000000, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Lts,
				regLoad(rs1, i, width64),
				regLoad(rs2, i, width64),
				expr.One,
				expr.Zero,
				width64,
			)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sltu",
		opcode:       opcode17(0b0000000, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Ltu,
				regLoad(rs1, i, width64),
				regLoad(rs2, i, width64),
				expr.One,
				expr.Zero,
				width64,
			)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "or",
		opcode:       opcode17(0b0000000, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Or, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "and",
		opcode:       opcode17(0b0000000, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.And, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "xor",
		opcode:       opcode17(0b0000000, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Xor, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sll",
		opcode:       opcode17(0b0000000, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := maskedRegOp(expr.Lsh, i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "srl",
		opcode:       opcode17(0b0000000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := maskedRegOp(expr.Rsh, i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sra",
		opcode:       opcode17(0b0100000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			shift := exprtools.MaskBits(regLoad(rs2, i, width64), 6, width64)
			val := exprtools.RshA(regLoad(rs1, i, width64), shift, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, { // Same definition as the 32bit version.
		name: "fence",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b0001111}),
			Mask:  revertBytes([]byte{0xf0, 0x0f, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeMemOrder,
		// FIXME: How to represent memory order?
		effects: func(i Instruction) []expr.Effect { return nil },
	}, { // Same definition as the 32bit version.
		name: "fence.i",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 1 << 4, 0b0001111}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeMemOrder,
		// FIXME: How to represent memory order?
		effects: func(i Instruction) []expr.Effect { return nil },
	}, { // Same definition as the 32bit version.
		name: "ecall",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
		// FIXME: How to represent syscall?
		effects: func(i Instruction) []expr.Effect { return nil },
	}, { // Same definition as the 32bit version.
		name: "ebreak",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 1 << 4, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
		// FIXME: How to represent syscall?
		effects: func(i Instruction) []expr.Effect { return nil },
	},

	// List of additional opcodes existing only in RV64.
	{
		name:         "addiw",
		opcode:       opcode10(0b000, 0b0011011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(regImmOp(expr.Add, immTypeI, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:                "slliw",
		opcode:              opcodeShiftImm(false, 5, 0b001, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(regImmShift(expr.Lsh, i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:                "srliw",
		opcode:              opcodeShiftImm(false, 5, 0b101, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(regImmShift(expr.Rsh, i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:                "sraiw",
		opcode:              opcodeShiftImm(true, 5, 0b101, 0b0011011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			reg := regLoad(rs1, i, width32)
			val := sext32To64(exprtools.RshA(reg, immShift(5, i), width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "addw",
		opcode:       opcode17(0b0000000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(reg2Op(expr.Add, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "subw",
		opcode:       opcode17(0b0100000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(reg2Op(expr.Sub, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sllw",
		opcode:       opcode17(0b0000000, 0b001, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(maskedRegOp(expr.Lsh, i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "srlw",
		opcode:       opcode17(0b0000000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(maskedRegOp(expr.Rsh, i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sraw",
		opcode:       opcode17(0b0100000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			shift := exprtools.MaskBits(regLoad(rs2, i, width32), 5, width32)
			val32 := exprtools.RshA(regLoad(rs1, i, width32), shift, width32)
			return []expr.Effect{regStore(sext32To64(val32), i, width64)}
		},
	},
}

var mul64 = []*instructionOpcode{
	{
		name:         "mul",
		opcode:       opcode17(0b0000001, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Mul, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "mulh",
		opcode:       opcode17(0b0000001, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			mul := exprtools.SignedMul(r1, r2, width64)
			shift := expr.NewConstUint[uint8](64, width8)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width128)
			val := exprtools.Crop(shifted, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "mulhu",
		opcode:       opcode17(0b0000001, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			mul := expr.NewBinary(expr.Mul, r1, r2, width128)
			shift := expr.NewConstUint[uint8](64, width8)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width128)
			val := exprtools.Crop(shifted, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "mulhsu",
		opcode:       opcode17(0b0000001, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			r1Abs := exprtools.Abs(r1, width64)
			mul := expr.NewBinary(expr.Mul, r1Abs, r2, width128)
			shift := expr.NewConstUint[uint8](64, width8)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width128)
			val := expr.NewCond(
				expr.Eq,
				r1,
				r1Abs,
				shifted,
				exprtools.Negate(shifted, width64),
				width64,
			)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "div",
		opcode:       opcode17(0b0000001, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			val := exprtools.SignedDiv(r1, r2, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "divu",
		opcode:       opcode17(0b0000001, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Div, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "rem",
		opcode:       opcode17(0b0000001, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			val := exprtools.SignedMod(r1, r2, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "remu",
		opcode:       opcode17(0b0000001, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Mod, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	},

	// rv64-only instructions:
	{
		name:         "mulw",
		opcode:       opcode17(0b0000001, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(reg2Op(expr.Mul, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "divw",
		opcode:       opcode17(0b0000001, 0b100, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			val := sext32To64(exprtools.SignedDiv(r1, r2, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "divuw",
		opcode:       opcode17(0b0000001, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(reg2Op(expr.Div, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "remw",
		opcode:       opcode17(0b0000001, 0b110, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			val := sext32To64(exprtools.SignedMod(r1, r2, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "remuw",
		opcode:       opcode17(0b0000001, 0b111, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		effects: func(i Instruction) []expr.Effect {
			val := sext32To64(reg2Op(expr.Mod, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	},
}
