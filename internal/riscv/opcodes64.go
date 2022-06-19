package riscv

import (
	"mltwist/internal/opcode"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"mltwist/pkg/model"
)

var integer64 = []*instructionType{
	{
		name:         "lui",
		opcode:       opcode7(0b0110111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		effects: func(i instruction) []expr.Effect {
			imm, _ := immTypeU.parseValue(i.value)
			val := sext32To64(expr.ConstFromInt(imm))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "auipc",
		opcode:       opcode7(0b0010111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeU,
		effects: func(i instruction) []expr.Effect {
			imm, _ := immTypeU.parseValue(i.value)
			val := expr.ConstFromUint(uint64(addrAddImm(i.addr, imm)))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "jal",
		opcode:       opcode7(0b1101111),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeJ,
		effects: func(i instruction) []expr.Effect {
			target := addrImmConst(immTypeJ, i, width64)
			// Address of following instruction.
			following := expr.ConstFromUint(uint64(i.addr) + 4)
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
		// FIXME: Find a way how to represent those jump targets.
		effects: func(i instruction) []expr.Effect {
			target := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
			// Address of following instruction.
			following := expr.ConstFromUint(uint64(i.addr) + 4)
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
		effects: func(i instruction) []expr.Effect {
			return []expr.Effect{branchCmp(exprtools.Eq, true, i, width64)}
		},
	}, {
		name:         "bne",
		opcode:       opcode10(0b001, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i instruction) []expr.Effect {
			return []expr.Effect{branchCmp(exprtools.Eq, false, i, width64)}
		},
	}, {
		name:         "blt",
		opcode:       opcode10(0b100, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i instruction) []expr.Effect {
			return []expr.Effect{branchCmp(exprtools.Lts, true, i, width64)}
		},
	}, {
		name:         "bge",
		opcode:       opcode10(0b101, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i instruction) []expr.Effect {
			return []expr.Effect{branchCmp(exprtools.Lts, false, i, width64)}
		},
	}, {
		name:         "bltu",
		opcode:       opcode10(0b110, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i instruction) []expr.Effect {
			return []expr.Effect{branchCmp(lessFunc, true, i, width64)}
		},
	}, {
		name:         "bgeu",
		opcode:       opcode10(0b111, 0b1100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		immediate:    immTypeB,
		effects: func(i instruction) []expr.Effect {
			return []expr.Effect{branchCmp(lessFunc, false, i, width64)}
		},
	}, {
		name:         "lb",
		opcode:       opcode10(0b000, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			addr := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
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
		effects: func(i instruction) []expr.Effect {
			addr := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
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
		effects: func(i instruction) []expr.Effect {
			addr := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
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
		effects: func(i instruction) []expr.Effect {
			addr := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
			val := memLoad(addr, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "lbu",
		opcode:       opcode10(0b100, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    1,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			addr := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
			return []expr.Effect{regStore(memLoad(addr, width8), i, width64)}
		},
	}, {
		name:         "lhu",
		opcode:       opcode10(0b101, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    2,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			addr := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
			return []expr.Effect{regStore(memLoad(addr, width16), i, width64)}
		},
	}, {
		name:         "lwu",
		opcode:       opcode10(0b110, 0b0000011),
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    4,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			addr := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
			return []expr.Effect{regStore(memLoad(addr, width32), i, width64)}
		},
	}, {
		name:         "sb",
		opcode:       opcode10(0b000, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   1,
		immediate:    immTypeS,
		effects: func(i instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(binOpFunc(expr.Add), immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width8)}
		},
	}, {
		name:         "sh",
		opcode:       opcode10(0b001, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   2,
		immediate:    immTypeS,
		effects: func(i instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(binOpFunc(expr.Add), immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width16)}
		},
	}, {
		name:         "sw",
		opcode:       opcode10(0b010, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   4,
		immediate:    immTypeS,
		effects: func(i instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(binOpFunc(expr.Add), immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width32)}
		},
	}, {
		name:         "sd",
		opcode:       opcode10(0b011, 0b0100011),
		inputRegCnt:  2,
		hasOutputReg: false,
		storeBytes:   8,
		immediate:    immTypeS,
		effects: func(i instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regImmOp(binOpFunc(expr.Add), immTypeS, i, width64)
			return []expr.Effect{memStore(val, addr, width64)}
		},
	}, {
		name:         "addi",
		opcode:       opcode10(0b000, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			val := regImmOp(binOpFunc(expr.Add), immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "slti",
		opcode:       opcode10(0b010, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			val := exprtools.Lts(
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
		effects: func(i instruction) []expr.Effect {
			val := expr.NewLess(
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
		effects: func(i instruction) []expr.Effect {
			val := regImmOp(exprtools.BitXor, immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "ori",
		opcode:       opcode10(0b110, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			val := regImmOp(exprtools.BitOr, immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "andi",
		opcode:       opcode10(0b111, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			val := regImmOp(exprtools.BitAnd, immTypeI, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "slli",
		opcode:       opcodeShiftImm(false, 6, 0b001, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := regImmShift(binOpFunc(expr.Lsh), i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "srli",
		opcode:       opcodeShiftImm(false, 6, 0b101, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := regImmShift(binOpFunc(expr.Rsh), i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "srai",
		opcode:       opcodeShiftImm(true, 6, 0b101, 0b0010011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := regImmShift(exprtools.RshA, i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "add",
		opcode:       opcode17(0b0000000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := reg2Op(binOpFunc(expr.Add), i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sub",
		opcode:       opcode17(0b0100000, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := reg2Op(exprtools.Sub, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "slt",
		opcode:       opcode17(0b0000000, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := exprtools.Lts(
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
		effects: func(i instruction) []expr.Effect {
			val := expr.NewLess(
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
		effects: func(i instruction) []expr.Effect {
			val := reg2Op(exprtools.BitOr, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "and",
		opcode:       opcode17(0b0000000, 0b111, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := reg2Op(exprtools.BitAnd, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "xor",
		opcode:       opcode17(0b0000000, 0b100, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := reg2Op(exprtools.BitXor, i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sll",
		opcode:       opcode17(0b0000000, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := maskedRegOp(binOpFunc(expr.Lsh), i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "srl",
		opcode:       opcode17(0b0000000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := maskedRegOp(binOpFunc(expr.Rsh), i, 6, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sra",
		opcode:       opcode17(0b0100000, 0b101, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := maskedRegOp(exprtools.RshA, i, 6, width64)
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
		effects: func(i instruction) []expr.Effect { return nil },
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
		effects: func(i instruction) []expr.Effect { return nil },
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
		effects: func(i instruction) []expr.Effect { return nil },
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
		effects: func(i instruction) []expr.Effect { return nil },
	},

	{
		name:         "csrrw",
		opcode:       opcode10(0b001, 0b1110011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i instruction) []expr.Effect {
			key := csrKey(i)
			return []expr.Effect{
				regStore(expr.NewRegLoad(key, width64), i, width64),
				expr.NewRegStore(regLoad(rs1, i, width64), key, width64),
			}
		},
	}, {
		name:         "csrrs",
		opcode:       opcode10(0b010, 0b1110011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width64)
			regVal := regLoad(rs1, i, width64)
			newVal := exprtools.BitOr(val, regVal, width64)
			return []expr.Effect{
				regStore(val, i, width64),
				expr.NewRegStore(newVal, key, width64),
			}
		},
	}, {
		name:         "csrrc",
		opcode:       opcode10(0b011, 0b1110011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width64)
			regVal := exprtools.BitNot(regLoad(rs1, i, width64), width64)
			newVal := exprtools.BitAnd(val, regVal, width64)
			return []expr.Effect{
				regStore(val, i, width64),
				expr.NewRegStore(newVal, key, width64),
			}
		},
	}, {
		name:         "csrrwi",
		opcode:       opcode10(0b101, 0b1110011),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i instruction) []expr.Effect {
			key := csrKey(i)
			return []expr.Effect{
				regStore(expr.NewRegLoad(key, width64), i, width64),
				expr.NewRegStore(csrImm(i), key, width64),
			}
		},
	}, {
		name:         "csrrsi",
		opcode:       opcode10(0b110, 0b1110011),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width64)
			newVal := exprtools.BitOr(val, csrImm(i), width64)
			return []expr.Effect{
				regStore(val, i, width64),
				expr.NewRegStore(newVal, key, width64),
			}
		},
	}, {
		name:         "csrrci",
		opcode:       opcode10(0b111, 0b1110011),
		inputRegCnt:  0,
		hasOutputReg: true,
		immediate:    immTypeI,
		instrType:    model.TypeCPUStateChange,
		effects: func(i instruction) []expr.Effect {
			key := csrKey(i)
			val := expr.NewRegLoad(key, width64)
			mask := exprtools.BitNot(csrImm(i), width64)
			newVal := exprtools.BitAnd(val, mask, width64)
			return []expr.Effect{
				regStore(val, i, width64),
				expr.NewRegStore(newVal, key, width64),
			}
		},
	},

	// List of additional opcodes existing only in RV64.
	{
		name:         "addiw",
		opcode:       opcode10(0b000, 0b0011011),
		inputRegCnt:  1,
		hasOutputReg: true,
		immediate:    immTypeI,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(regImmOp(binOpFunc(expr.Add), immTypeI, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "slliw",
		opcode:       opcodeShiftImm(false, 5, 0b001, 0b0011011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(regImmShift(binOpFunc(expr.Lsh), i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "srliw",
		opcode:       opcodeShiftImm(false, 5, 0b101, 0b0011011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(regImmShift(binOpFunc(expr.Rsh), i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sraiw",
		opcode:       opcodeShiftImm(true, 5, 0b101, 0b0011011),
		inputRegCnt:  1,
		hasOutputReg: true,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(regImmShift(exprtools.RshA, i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "addw",
		opcode:       opcode17(0b0000000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(reg2Op(binOpFunc(expr.Add), i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "subw",
		opcode:       opcode17(0b0100000, 0b000, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(reg2Op(exprtools.Sub, i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sllw",
		opcode:       opcode17(0b0000000, 0b001, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(maskedRegOp(binOpFunc(expr.Lsh), i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "srlw",
		opcode:       opcode17(0b0000000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(maskedRegOp(binOpFunc(expr.Rsh), i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sraw",
		opcode:       opcode17(0b0100000, 0b101, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(maskedRegOp(exprtools.RshA, i, 5, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	},
}

var mul64 = []*instructionType{
	{
		name:         "mul",
		opcode:       opcode17(0b0000001, 0b000, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			val := reg2Op(binOpFunc(expr.Mul), i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "mulh",
		opcode:       opcode17(0b0000001, 0b001, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			mul := exprtools.SignedMul(r1, r2, width64)
			shift := expr.ConstFromUint[uint8](64)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width128)
			val := exprtools.NewWidthGadget(shifted, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "mulhu",
		opcode:       opcode17(0b0000001, 0b011, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			mul := expr.NewBinary(expr.Mul, r1, r2, width128)
			shift := expr.ConstFromUint[uint8](64)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width128)
			val := exprtools.NewWidthGadget(shifted, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "mulhsu",
		opcode:       opcode17(0b0000001, 0b010, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			r1Abs := exprtools.Abs(r1, width64)
			mul := expr.NewBinary(expr.Mul, r1Abs, r2, width128)
			shift := expr.ConstFromUint[uint8](64)
			shifted := expr.NewBinary(expr.Rsh, mul, shift, width128)
			val := exprtools.BoolCond(
				exprtools.IntNegative(r1, width64),
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
		effects: func(i instruction) []expr.Effect {
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
		effects: func(i instruction) []expr.Effect {
			val := reg2Op(binOpFunc(expr.Div), i, width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "rem",
		opcode:       opcode17(0b0000001, 0b110, 0b0110011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
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
		effects: func(i instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width64), regLoad(rs2, i, width64)
			val := exprtools.Mod(r1, r2, width64)
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
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(reg2Op(binOpFunc(expr.Mul), i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "divw",
		opcode:       opcode17(0b0000001, 0b100, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
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
		effects: func(i instruction) []expr.Effect {
			val := sext32To64(reg2Op(binOpFunc(expr.Div), i, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "remw",
		opcode:       opcode17(0b0000001, 0b110, 0b0111011),
		inputRegCnt:  2,
		hasOutputReg: true,
		immediate:    immTypeR,
		effects: func(i instruction) []expr.Effect {
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
		effects: func(i instruction) []expr.Effect {
			r1, r2 := regLoad(rs1, i, width32), regLoad(rs2, i, width32)
			val := sext32To64(exprtools.Mod(r1, r2, width32))
			return []expr.Effect{regStore(val, i, width64)}
		},
	},
}

var atomic64 = []*instructionType{
	{
		name: "lr.d",
		// LR.D is special as rs2 is always r0 - otherwise the
		// instruction opcode is undefined as zeros are defined in an
		// instruction encoding.
		opcode: opcode.Opcode{
			Bytes: []byte{0b0101111, 0b011 << 4, 0, 0b00010 << 3},
			Mask:  []byte{0x7f, 0b111 << 4, 0xf0, 0xf9},
		},
		inputRegCnt:  1,
		hasOutputReg: true,
		loadBytes:    8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			val := memLoad(regLoad(rs1, i, width64), width64)
			return []expr.Effect{regStore(val, i, width64)}
		},
	}, {
		name:         "sc.d",
		opcode:       opcodeAtomic(0b00011, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			val := regLoad(rs2, i, width64)
			addr := regLoad(rs1, i, width64)
			// TODO: Find a way how to emulate the race check.
			return []expr.Effect{
				memStore(val, addr, width64),
				regStore(expr.Zero, i, width64),
			}
		},
	}, {
		name:         "amoswap.d",
		opcode:       opcodeAtomic(0b00001, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			addr := regLoad(rs1, i, width64)
			return []expr.Effect{
				regStore(memLoad(addr, width64), i, width64),
				memStore(regLoad(rs2, i, width64), addr, width64),
			}
		},
	}, {
		name:         "amoadd.d",
		opcode:       opcodeAtomic(0, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(binOpFunc(expr.Add), i, width64)
		},
	}, {
		name:         "amoxor.d",
		opcode:       opcodeAtomic(0b00100, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(exprtools.BitXor, i, width64)
		},
	}, {
		name:         "amoand.d",
		opcode:       opcodeAtomic(0b01100, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(exprtools.BitAnd, i, width64)
		},
	}, {
		name:         "amoor.d",
		opcode:       opcodeAtomic(0b01000, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(exprtools.BitOr, i, width64)
		},
	}, {
		name:         "amomin.d",
		opcode:       opcodeAtomic(0b10000, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(atomicMinMax(exprtools.Lts, false), i, width64)
		},
	}, {
		name:         "amomax.d",
		opcode:       opcodeAtomic(0b10100, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(atomicMinMax(exprtools.Lts, true), i, width64)
		},
	}, {
		name:         "amominu.d",
		opcode:       opcodeAtomic(0b11000, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(atomicMinMax(lessFunc, false), i, width64)
		},
	}, {
		name:         "amomaxu.d",
		opcode:       opcodeAtomic(0b11100, 0b011, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    8,
		storeBytes:   8,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOp(atomicMinMax(lessFunc, true), i, width64)
		},
	},

	// RV64 specific instructions operating on 32 bit values.
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
		effects: func(i instruction) []expr.Effect {
			val := memLoad(regLoad(rs1, i, width64), width32)
			return []expr.Effect{regStore(sext32To64(val), i, width64)}
		},
	}, {
		name:         "sc.w",
		opcode:       opcodeAtomic(0b00011, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			val := regLoad(rs2, i, width32)
			addr := regLoad(rs1, i, width64)
			// TODO: Find a way how to emulate the race check.
			return []expr.Effect{
				memStore(val, addr, width32),
				regStore(expr.Zero, i, width64),
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
		effects: func(i instruction) []expr.Effect {
			addr := regLoad(rs1, i, width64)
			return []expr.Effect{
				regStore(sext32To64(memLoad(addr, width32)), i, width64),
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
		effects: func(i instruction) []expr.Effect {
			return atomicOpWidth(binOpFunc(expr.Add), i, width64, width32)
		},
	}, {
		name:         "amoxor.w",
		opcode:       opcodeAtomic(0b00100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOpWidth(exprtools.BitXor, i, width64, width32)
		},
	}, {
		name:         "amoand.w",
		opcode:       opcodeAtomic(0b01100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOpWidth(exprtools.BitAnd, i, width64, width32)
		},
	}, {
		name:         "amoor.w",
		opcode:       opcodeAtomic(0b01000, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			return atomicOpWidth(exprtools.BitOr, i, width64, width32)
		},
	}, {
		name:         "amomin.w",
		opcode:       opcodeAtomic(0b10000, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			f := atomicMinMax(exprtools.Lts, false)
			return atomicOpWidth(f, i, width64, width32)
		},
	}, {
		name:         "amomax.w",
		opcode:       opcodeAtomic(0b10100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			f := atomicMinMax(exprtools.Lts, true)
			return atomicOpWidth(f, i, width64, width32)
		},
	}, {
		name:         "amominu.w",
		opcode:       opcodeAtomic(0b11000, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			f := atomicMinMax(lessFunc, false)
			return atomicOpWidth(f, i, width64, width32)
		},
	}, {
		name:         "amomaxu.w",
		opcode:       opcodeAtomic(0b11100, 0b010, 0b0101111),
		inputRegCnt:  2,
		hasOutputReg: true,
		loadBytes:    4,
		storeBytes:   4,
		instrType:    model.TypeMemOrder,
		effects: func(i instruction) []expr.Effect {
			f := atomicMinMax(lessFunc, true)
			return atomicOpWidth(f, i, width64, width32)
		},
	},
}
