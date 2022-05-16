package riscv

import (
	"mltwist/internal/opcode"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"mltwist/pkg/model"
)

var integer32 = []*instructionOpcode{
	{
		name:         "lui",
		opcode:       opcode7(0b0110111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		effects: func(i Instruction) []expr.Effect {
			imm, _ := immTypeU.parseValue(i.value)
			val := expr.ConstFromInt(imm)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "auipc",
		opcode:       opcode7(0b0010111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		effects: func(i Instruction) []expr.Effect {
			imm, _ := immTypeU.parseValue(i.value)
			val := expr.ConstFromUint(uint32(addrAddImm(i.address, imm)))
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "jal",
		opcode:       opcode7(0b1101111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeJ,
		effects: func(i Instruction) []expr.Effect {
			target := addrImmConst(immTypeJ, i, width32)
			// Address of following instruction.
			following := expr.ConstFromUint(uint32(i.address) + 4)
			return []expr.Effect{
				expr.NewRegStore(target, expr.IPKey, width32),
				regStore(following, i, width32),
			}
		},
	}, {
		name:         "jalr",
		opcode:       opcode10(0b000, 0b1100111),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		// FIXME: Find a way how to represent those jump targets.
		effects: func(i Instruction) []expr.Effect {
			target := regImmOp(expr.Add, immTypeI, i, width32)
			// Address of following instruction.
			following := expr.ConstFromUint(uint32(i.address) + 4)
			return []expr.Effect{
				expr.NewRegStore(target, expr.IPKey, width32),
				regStore(following, i, width32),
			}
		},
	}, {
		name:         "beq",
		opcode:       opcode10(0b000, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Eq, true, i, width32)}
		},
	}, {
		name:         "bne",
		opcode:       opcode10(0b001, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Eq, false, i, width32)}
		},
	}, {
		name:         "blt",
		opcode:       opcode10(0b100, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Lts, true, i, width32)}
		},
	}, {
		name:         "bge",
		opcode:       opcode10(0b101, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Lts, false, i, width32)}
		},
	}, {
		name:         "bltu",
		opcode:       opcode10(0b110, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Ltu, true, i, width32)}
		},
	}, {
		name:         "bgeu",
		opcode:       opcode10(0b111, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i Instruction) []expr.Effect {
			return []expr.Effect{branchCmp(expr.Ltu, false, i, width32)}
		},
	}, {
		name:         "lb",
		opcode:       opcode10(0b000, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			val := sext(memLoad(addr, width8), 7, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "lh",
		opcode:       opcode10(0b001, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			val := sext(memLoad(addr, width16), 15, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "lw",
		opcode:       opcode10(0b010, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			val := memLoad(addr, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "lbu",
		opcode:       opcode10(0b100, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			return []expr.Effect{regStore(memLoad(addr, width8), i, width32)}
		},
	}, {
		name:         "lhu",
		opcode:       opcode10(0b101, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			addr := regImmOp(expr.Add, immTypeI, i, width32)
			return []expr.Effect{regStore(memLoad(addr, width16), i, width32)}
		},
	}, {
		name:         "sb",
		opcode:       opcode10(0b000, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   1,
		immediate:    immTypeS,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width32)
			addr := regImmOp(expr.Add, immTypeS, i, width32)
			return []expr.Effect{memStore(val, addr, width8)}
		},
	}, {
		name:         "sh",
		opcode:       opcode10(0b001, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   2,
		immediate:    immTypeS,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width32)
			addr := regImmOp(expr.Add, immTypeS, i, width32)
			return []expr.Effect{memStore(val, addr, width16)}
		},
	}, {
		name:         "sw",
		opcode:       opcode10(0b010, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   4,
		immediate:    immTypeS,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width32)
			addr := regImmOp(expr.Add, immTypeS, i, width32)
			return []expr.Effect{memStore(val, addr, width32)}
		},
	}, {
		name:         "addi",
		opcode:       opcode10(0b000, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.Add, immTypeI, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "slti",
		opcode:       opcode10(0b010, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Lts,
				regLoad(rs1, i, width32),
				immConst(immTypeI, i),
				expr.One,
				expr.Zero,
				width32,
			)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "sltiu",
		opcode:       opcode10(0b011, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Ltu,
				regLoad(rs1, i, width32),
				immConst(immTypeI, i),
				expr.One,
				expr.Zero,
				width32,
			)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "xori",
		opcode:       opcode10(0b100, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.Xor, immTypeI, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "ori",
		opcode:       opcode10(0b110, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.Or, immTypeI, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "andi",
		opcode:       opcode10(0b111, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i Instruction) []expr.Effect {
			val := regImmOp(expr.And, immTypeI, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "slli",
		opcode:       opcodeShiftImm(false, 5, 0b001, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := regImmShift(expr.Lsh, i, 5, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "srli",
		opcode:       opcodeShiftImm(false, 5, 0b101, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := regImmShift(expr.Rsh, i, 5, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "srai",
		opcode:       opcodeShiftImm(true, 5, 0b101, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			reg := regLoad(rs1, i, width32)
			val := exprtools.RshA(reg, immShift(5, i), width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "add",
		opcode:       opcode17(0b0000000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Add, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "sub",
		opcode:       opcode17(0b0100000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Sub, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "slt",
		opcode:       opcode17(0b0000000, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Lts,
				regLoad(rs1, i, width32),
				regLoad(rs2, i, width32),
				expr.One,
				expr.Zero,
				width32,
			)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "sltu",
		opcode:       opcode17(0b0000000, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := expr.NewCond(
				expr.Ltu,
				regLoad(rs1, i, width32),
				regLoad(rs2, i, width32),
				expr.One,
				expr.Zero,
				width32,
			)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "or",
		opcode:       opcode17(0b0000000, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Or, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "and",
		opcode:       opcode17(0b0000000, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.And, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "xor",
		opcode:       opcode17(0b0000000, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Xor, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "sll",
		opcode:       opcode17(0b0000000, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := maskedRegOp(expr.Lsh, i, 5, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "srl",
		opcode:       opcode17(0b0000000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			val := maskedRegOp(expr.Rsh, i, 5, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "sra",
		opcode:       opcode17(0b0100000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i Instruction) []expr.Effect {
			shift := exprtools.MaskBits(regLoad(rs2, i, width32), 5, width32)
			val := exprtools.RshA(regLoad(rs1, i, width32), shift, width32)
			return []expr.Effect{regStore(val, i, width32)}
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
		effects:      func(i Instruction) []expr.Effect { return nil },
	}, {
		name: "fence.i",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 1 << 4, 0b0001111}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeMemOrder,
		effects:      func(i Instruction) []expr.Effect { return nil },
	}, {
		name: "ecall",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 0, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
		effects:      func(i Instruction) []expr.Effect { return nil },
	}, {
		name: "ebreak",
		opcode: opcode.Opcode{
			Bytes: revertBytes([]byte{0, 1 << 4, 0, 0b1110011}),
			Mask:  revertBytes([]byte{0xff, 0xff, 0xff, 0xff}),
		},
		inputRegCnt:  0,
		hasOutputReg: false,
		instrType:    model.TypeSyscall,
		effects:      func(i Instruction) []expr.Effect { return nil },
	},

	// TODO: Find a way how to represent CSR instructions.

	{
		name:         "csrrw",
		opcode:       opcode10(0b001, 0b1110011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i Instruction) []expr.Effect {
			key := csrKey(i)
			return []expr.Effect{
				regStore(expr.NewRegLoad(key, width32), i, width32),
				expr.NewRegStore(regLoad(rs1, i, width32), key, width32),
			}
		},
	}, {
		name:         "csrrs",
		opcode:       opcode10(0b010, 0b1110011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i Instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width32)
			mask := regLoad(rs1, i, width32)
			newVal := expr.NewBinary(expr.Or, val, mask, width32)
			return []expr.Effect{
				regStore(val, i, width32),
				expr.NewRegStore(newVal, key, width32),
			}
		},
	}, {
		name:         "csrrc",
		opcode:       opcode10(0b011, 0b1110011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i Instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width32)
			mask := exprtools.BitNegate(regLoad(rs1, i, width32), width32)
			newVal := expr.NewBinary(expr.And, val, mask, width32)
			return []expr.Effect{
				regStore(val, i, width32),
				expr.NewRegStore(newVal, key, width32),
			}
		},
	}, {
		name:         "csrrwi",
		opcode:       opcode10(0b101, 0b1110011),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i Instruction) []expr.Effect {
			key := csrKey(i)
			return []expr.Effect{
				regStore(expr.NewRegLoad(key, width32), i, width32),
				expr.NewRegStore(csrImm(i), key, width32),
			}
		},
	}, {
		name:         "csrrsi",
		opcode:       opcode10(0b110, 0b1110011),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i Instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width32)
			newVal := expr.NewBinary(expr.Or, val, csrImm(i), width32)
			return []expr.Effect{
				regStore(val, i, width32),
				expr.NewRegStore(newVal, key, width32),
			}
		},
	}, {
		name:         "csrrci",
		opcode:       opcode10(0b111, 0b1110011),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i Instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width32)
			mask := exprtools.BitNegate(csrImm(i), width32)
			newVal := expr.NewBinary(expr.And, val, mask, width32)
			return []expr.Effect{
				regStore(val, i, width32),
				expr.NewRegStore(newVal, key, width32),
			}
		},
	},
}

var mul32 = []*instructionOpcode{
	{
		name:         "mul",
		opcode:       opcode17(0b0000001, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Mul, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "mulh",
		opcode:       opcode17(0b0000001, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			mul := exprtools.SignedMul(r1, r2, width32)
			shift := expr.ConstFromUint[uint8](32)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width64)
			val := exprtools.NewWidthGadget(shifted, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "mulhu",
		opcode:       opcode17(0b0000001, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			mul := expr.NewBinary(expr.Mul, r1, r2, width64)
			shift := expr.ConstFromUint[uint8](32)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width64)
			val := exprtools.NewWidthGadget(shifted, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "mulhsu",
		opcode:       opcode17(0b0000001, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			r1Abs := exprtools.Abs(r1, width32)
			mul := expr.NewBinary(expr.Mul, r1Abs, r2, width64)
			shift := expr.ConstFromUint[uint8](32)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width64)
			val := expr.NewCond(
				expr.Eq,
				r1,
				r1Abs,
				shifted,
				exprtools.Negate(shifted, width32),
				width32,
			)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "div",
		opcode:       opcode17(0b0000001, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			val := exprtools.SignedDiv(r1, r2, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "divu",
		opcode:       opcode17(0b0000001, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Div, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "rem",
		opcode:       opcode17(0b0000001, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			val := exprtools.SignedMod(r1, r2, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "remu",
		opcode:       opcode17(0b0000001, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i Instruction) []expr.Effect {
			val := reg2Op(expr.Mod, i, width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	},
}

var atomic32 = []*instructionOpcode{
	{
		name: "lr.w",
		// LR.W is special as rs2 is always r0 - otherwise the
		// instruction opcode is undefined as zeros are defined in an
		// instruction encoding.
		opcode: opcode.Opcode{
			Bytes: []byte{0b0101111, 0b010 << 4, 0, 0b00010 << 3},
			Mask:  []byte{0x7f, 0b111 << 4, 0xf0, 0xf9},
		},
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			val := memLoad(regLoad(rs1, i, width32), width32)
			return []expr.Effect{regStore(val, i, width32)}
		},
	}, {
		name:         "sc.w",
		opcode:       opcodeAtomic(0b00011, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			val := regLoad(rs2, i, width32)
			addr := regLoad(rs1, i, width32)
			// TODO: Find a way how to emulate the race check.
			return []expr.Effect{
				memStore(val, addr, width32),
				regStore(expr.Zero, i, width32),
			}
		},
	}, {
		name:         "amoswap.w",
		opcode:       opcodeAtomic(0b00001, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			addr := regLoad(rs1, i, width32)
			return []expr.Effect{
				regStore(memLoad(addr, width32), i, width32),
				memStore(regLoad(rs2, i, width32), addr, width32),
			}
		},
	}, {
		name:         "amoadd.W",
		opcode:       opcodeAtomic(0, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicBinaryOp(expr.Add), i, width32)
		},
	}, {
		name:         "amoxor.w",
		opcode:       opcodeAtomic(0b00100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicBinaryOp(expr.Xor), i, width32)
		},
	}, {
		name:         "amoand.w",
		opcode:       opcodeAtomic(0b01100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicBinaryOp(expr.And), i, width32)
		},
	}, {
		name:         "amoor.w",
		opcode:       opcodeAtomic(0b01000, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicBinaryOp(expr.Or), i, width32)
		},
	}, {
		name:         "amomin.w",
		opcode:       opcodeAtomic(0b10000, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicMinMax(expr.Lts, false), i, width32)
		},
	}, {
		name:         "amomax.w",
		opcode:       opcodeAtomic(0b10100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicMinMax(expr.Lts, true), i, width32)
		},
	}, {
		name:         "amominu.w",
		opcode:       opcodeAtomic(0b11000, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicMinMax(expr.Ltu, false), i, width32)
		},
	}, {
		name:         "amomaxu.w",
		opcode:       opcodeAtomic(0b11100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i Instruction) []expr.Effect {
			return atomicOp(atomicMinMax(expr.Ltu, true), i, width32)
		},
	},
}
