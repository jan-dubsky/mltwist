package riscv

import (
	"decomp/internal/opcode"
	"decomp/pkg/expr"
	"decomp/pkg/expr/exprtools"
	"decomp/pkg/model"
)

var integer32 = []*instructionOpcode{
	{
		name:         "lui",
		opcode:       opcode7(0b0110111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			imm, _ := immTypeU.parseValue(i.value)
			return expr.NewConstInt(imm, width32)
		},
	}, {
		name:         "auipc",
		opcode:       opcode7(0b0010111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			imm, _ := immTypeU.parseValue(i.value)
			return expr.NewConstUint(addrAddImm(i.address, imm), width32)
		},
	}, {
		name:         "jal",
		opcode:       opcode7(0b1101111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeJ,
		instrType:    model.TypeJump,
		jumpTarget: func(i Instruction) model.Address {
			imm, _ := immTypeJ.parseValue(i.value)
			return addrAddImm(i.address, imm)
		},
		expr: func(i Instruction) expr.Expr {
			// Address of following instruction.
			return expr.NewConstUint(uint64(i.address+4), width32)
		},
	}, {
		name:         "jalr",
		opcode:       opcode10(0b000, 0b1100111),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeJumpDyn,
		// FIXME: Find a way how to represent those jump targets.
		expr: func(i Instruction) expr.Expr {
			// Address of following instruction.
			return expr.NewConstUint(uint64(i.address+4), width32)
		},
	}, {
		name:         "beq",
		opcode:       opcode10(0b000, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		jumpTarget:   branchJumpTarget,
	}, {
		name:         "bne",
		opcode:       opcode10(0b001, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		jumpTarget:   branchJumpTarget,
	}, {
		name:         "blt",
		opcode:       opcode10(0b100, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		jumpTarget:   branchJumpTarget,
	}, {
		name:         "bge",
		opcode:       opcode10(0b101, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		jumpTarget:   branchJumpTarget,
	}, {
		name:         "bltu",
		opcode:       opcode10(0b110, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		jumpTarget:   branchJumpTarget,
	}, {
		name:         "bgeu",
		opcode:       opcode10(0b111, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		instrType:    model.TypeCJump,
		jumpTarget:   branchJumpTarget,
	}, {
		name:         "lb",
		opcode:       opcode10(0b000, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			return sext(expr.NewLoad(addr, width8), 7, width32)
		},
	}, {
		name:         "lh",
		opcode:       opcode10(0b001, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			return sext(expr.NewLoad(addr, width16), 15, width32)
		},
	}, {
		name:         "lw",
		opcode:       opcode10(0b010, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		immediate:    immTypeI,
		instrType:    model.TypeLoad,
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			return expr.NewLoad(addr, width32)
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
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			return exprtools.Crop(expr.NewLoad(addr, width8), width32)
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
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			return exprtools.Crop(expr.NewLoad(addr, width16), width32)
		},
	}, {
		name:         "sb",
		opcode:       opcode10(0b000, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   1,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeS, i, width32)
			return expr.NewStore(regExpr(rs2, i, width32), addr, width8)
		},
	}, {
		name:         "sh",
		opcode:       opcode10(0b001, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   2,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeS, i, width32)
			return expr.NewStore(regExpr(rs2, i, width32), addr, width16)
		},
	}, {
		name:         "sw",
		opcode:       opcode10(0b010, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   4,
		immediate:    immTypeS,
		instrType:    model.TypeStore,
		expr: func(i Instruction) expr.Expr {
			addr := regImmOp(expr.Add, immTypeS, i, width32)
			return expr.NewStore(regExpr(rs2, i, width32), addr, width32)
		},
	}, {
		name:         "addi",
		opcode:       opcode10(0b000, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return regImmOp(expr.Add, immTypeI, i, width32)
		},
	}, {
		name:         "slti",
		opcode:       opcode10(0b010, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return expr.NewCond(
				expr.Lt,
				regExpr(rs1, i, width32),
				immConst(immTypeI, i),
				expr.One,
				expr.Zero,
				width32,
			)
		},
	}, {
		name:         "sltiu",
		opcode:       opcode10(0b011, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return expr.NewCond(
				expr.Ltu,
				regExpr(rs1, i, width32),
				immConst(immTypeI, i),
				expr.One,
				expr.Zero,
				width32,
			)
		},
	}, {
		name:         "xori",
		opcode:       opcode10(0b100, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return regImmOp(expr.Xor, immTypeI, i, width32)
		},
	}, {
		name:         "ori",
		opcode:       opcode10(0b110, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return regImmOp(expr.Or, immTypeI, i, width32)
		},
	}, {
		name:         "andi",
		opcode:       opcode10(0b111, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return regImmOp(expr.And, immTypeI, i, width32)
		},
	}, {
		name:                "slli",
		opcode:              opcodeShiftImm(false, 5, 0b001, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return regImmShift(expr.Lsh, i, 5, width32)
		},
	}, {
		name:                "srli",
		opcode:              opcodeShiftImm(false, 5, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return regImmShift(expr.Rsh, i, 5, width32)
		},
	}, {
		name:                "srai",
		opcode:              opcodeShiftImm(true, 5, 0b101, 0b0010011),
		inputRegCnt:         1,
		hasOutputReg:        true,
		additionalImmediate: addImmSh32,
		instrType:           model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return regImmShift(expr.RshL, i, 5, width32)
		},
	}, {
		name:         "add",
		opcode:       opcode17(0b0000000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.Add, i, width32)
		},
	}, {
		name:         "sub",
		opcode:       opcode17(0b0100000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.Sub, i, width32)
		},
	}, {
		name:         "slt",
		opcode:       opcode17(0b0000000, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return expr.NewCond(
				expr.Lt,
				regExpr(rs1, i, width32),
				regExpr(rs2, i, width32),
				expr.One,
				expr.Zero,
				width32,
			)
		},
	}, {
		name:         "sltu",
		opcode:       opcode17(0b0000000, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return expr.NewCond(
				expr.Ltu,
				regExpr(rs1, i, width32),
				regExpr(rs2, i, width32),
				expr.One,
				expr.Zero,
				width32,
			)
		},
	}, {
		name:         "or",
		opcode:       opcode17(0b0000000, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.Or, i, width32)
		},
	}, {
		name:         "and",
		opcode:       opcode17(0b0000000, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.And, i, width32)
		},
	}, {
		name:         "xor",
		opcode:       opcode17(0b0000000, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.Xor, i, width32)
		},
	}, {
		name:         "sll",
		opcode:       opcode17(0b0000000, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return maskedRegOp(expr.Lsh, i, 5, width32)
		},
	}, {
		name:         "srl",
		opcode:       opcode17(0b0000000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return maskedRegOp(expr.Rsh, i, 5, width32)
		},
	}, {
		name:         "sra",
		opcode:       opcode17(0b0100000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return maskedRegOp(expr.RshL, i, 5, width32)
		},
	}, {
		name: "fence",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b0001111}),
			Mask:  revertBytes([]byte{0xf0, 0x0f, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeMemOrder,
	}, {
		name: "fence.i",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 1 << 4, 0b0001111}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeMemOrder,
	}, {
		name: "ecall",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
	}, {
		name: "ebreak",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 1 << 4, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
	},

	// TODO: Find a way how to represent CSR instructions.
	/*
		{
			name:                "csrrw",
			opcode:              opcode10(0b001, 0b1110011),
			inputRegCnt:         1,
			hasOutputReg:        true,
			immediate:           immTypeI,
			additionalImmediate: addImmCSR,
			instrType:           model.TypeCPUStateChange,
		}, {
			name:                "csrrs",
			opcode:              opcode10(0b010, 0b1110011),
			inputRegCnt:         1,
			hasOutputReg:        true,
			immediate:           immTypeI,
			additionalImmediate: addImmCSR,
			instrType:           model.TypeCPUStateChange,
		}, {
			name:                "csrrc",
			opcode:              opcode10(0b011, 0b1110011),
			inputRegCnt:         1,
			hasOutputReg:        true,
			immediate:           immTypeI,
			additionalImmediate: addImmCSR,
			instrType:           model.TypeCPUStateChange,
		},
		// FIXME: There are 2 independent immediate values in those 3
		// instructions. Find a way how to parse and represent those
		// instructions.
		{
			name:                "csrrwi",
			opcode:              opcode10(0b101, 0b1110011),
			inputRegCnt:         0,
			hasOutputReg:        true,
			additionalImmediate: addImmCSR,
			instrType:           model.TypeCPUStateChange,
		}, {
			name:                "csrrsi",
			opcode:              opcode10(0b110, 0b1110011),
			inputRegCnt:         0,
			hasOutputReg:        true,
			additionalImmediate: addImmCSR,
			instrType:           model.TypeCPUStateChange,
		}, {
			name:                "csrrci",
			opcode:              opcode10(0b111, 0b1110011),
			inputRegCnt:         0,
			hasOutputReg:        true,
			additionalImmediate: addImmCSR,
			instrType:           model.TypeCPUStateChange,
		},
	*/
}

var mul32 = []*instructionOpcode{
	{
		name:         "mul",
		opcode:       opcode17(0b0000001, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.Mul, i, width32)
		},
	}, {
		name:         "mulh",
		opcode:       opcode17(0b0000001, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			r1, r2 := regExpr(rs1, i, width32), regExpr(rs2, i, width32)
			mul := exprtools.SignedMul(r1, r2, width32)
			shift := expr.NewConstUint(uint8(32), width8)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width64)
			return exprtools.Crop(shifted, width32)
		},
	}, {
		name:         "mulhu",
		opcode:       opcode17(0b0000001, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			r1, r2 := regExpr(rs1, i, width32), regExpr(rs2, i, width32)
			mul := expr.NewBinary(expr.Mul, r1, r2, width64)
			shift := expr.NewConstUint(uint8(32), width8)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width64)
			return exprtools.Crop(shifted, width32)
		},
	}, {
		name:         "mulhsu",
		opcode:       opcode17(0b0000001, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			r1, r2 := regExpr(rs1, i, width32), regExpr(rs2, i, width32)
			r1Abs := exprtools.Abs(r1, width32)
			mul := expr.NewBinary(expr.Mul, r1Abs, r2, width64)
			shift := expr.NewConstUint(uint8(32), width8)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width64)
			return expr.NewCond(
				expr.Eq,
				r1,
				r1Abs,
				shifted,
				exprtools.Negate(shifted, width32),
				width32,
			)
		},
	}, {
		name:         "div",
		opcode:       opcode17(0b0000001, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			r1, r2 := regExpr(rs1, i, width32), regExpr(rs2, i, width32)
			return exprtools.SignedDiv(r1, r2, width32)
		},
	}, {
		name:         "divu",
		opcode:       opcode17(0b0000001, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.Div, i, width32)
		},
	}, {
		name:         "rem",
		opcode:       opcode17(0b0000001, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			r1, r2 := regExpr(rs1, i, width32), regExpr(rs2, i, width32)
			return exprtools.SignedMod(r1, r2, width32)
		},
	}, {
		name:         "remu",
		opcode:       opcode17(0b0000001, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		instrType:    model.TypeAritm,
		expr: func(i Instruction) expr.Expr {
			return reg2Op(expr.Mod, i, width32)
		},
	},
}
