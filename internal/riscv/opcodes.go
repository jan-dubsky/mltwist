package riscv

import (
	"fmt"
	"mltwist/internal/opcode"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"mltwist/pkg/model"
)

// MemoryKey is identifier of the memory address space.
//
// Given that RISCV specification defines only one memory space, there is no
// reason to explain which memory is identified by this key.
const MemoryKey = expr.Key("memory")

const (
	// low7Bits is a byte (mask) with bottom 7 bits set and last bit unset.
	low7Bits byte = 0x7F
	// low3Bits is a byte (mask) with bottom 3 bits set and all higher bits
	// unset.
	low3Bits byte = 0x7
	// low5Bits is a byte (mask) with bottom 5 bits set and all higher bits
	// unset.
	low5Bits byte = 0x1f
)

// assertMask checks that only bits set in mask are set in b. This method will
// panic if any other bit is set on b.
func assertMask(b byte, mask byte) {
	if b&mask != b {
		panic(fmt.Sprintf("bits must match mask 0x%x: 0x%x", mask, b))
	}
}

// revertBytes is an utility function which reverts b as a slice. For more
// convenience, this function also returns b.
//
// The purpose of this method is to lower cognitive complexity of RISC-V
// instruction definition in the code. The problem is that even though RISC-V
// instructions are encoded in little endian byte order, the RICS-V
// specification ifself user big endian notation (least significant byte is the
// right-most byte). This misalignment in between ste specs and real
// implementation is expected. In human readable documents, it's just common to
// write least significant byte and bit to right. But as humans are not good in
// reverting byte sequences they read, it's much better solution to write
// instruction opcodes in big endian and then to revert the array.
func revertBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// opcode7 returns opcode matching low with bottom 7 bits of an instruction.
//
// As 1 byte opcodes have only 7 bits, this method will panic for values of low
// greater than 127.
func opcode7(low byte) opcode.Opcode {
	assertMask(low, low7Bits)

	return opcode.Opcode{
		Bytes: []byte{low},
		Mask:  []byte{low7Bits},
	}
}

// opcode10 returns opcode matching low with bottom 7 bits and mid with bits
// [12..14] of an instruction.
//
// This method panics if low is greater than 127 or if mid is greater than 7.
func opcode10(mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)

	return opcode.Opcode{
		Bytes: []byte{low, mid << 4},
		Mask:  []byte{low7Bits, low3Bits << 4},
	}
}

// opcode10 returns opcode matching low with bottom 7 bits, mid with bits
// [12..14] and high with bits [25..31] of an instruction.
//
// This method panics if low is greater than 127, mid is greater than 7 or if
// high is greater than 127.
func opcode17(high byte, mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)
	assertMask(high, low7Bits)

	return opcode.Opcode{
		Bytes: []byte{low, mid << 4, 0, high << 1},
		Mask:  []byte{low7Bits, low3Bits << 4, 0, low7Bits << 1},
	}
}

func assertShiftBits(shiftBits uint8) {
	if s := shiftBits; s != 5 && s != 6 {
		panic(fmt.Sprintf("invalid immediate-encoded shift bit count: %d", s))
	}
}

// opcodeShiftImm creates an opcode definition for RISC V bit shift instruction
// with shift immediate encoded in an instruction opcode.
//
// Even though RISC V manual states that there are only 6 distinct instruction
// encodings and all of them should be describable by either opcode7, opcode10,
// or opcode17, there is one small exception. Yes x86, you are not the only one
// architecture doing weird things... The exceptional opcode encoding is the one
// with fixed bit short argument in the instruction immediate value.
// Technically, such an instruction is just a bit shift with 12 bit immediate,
// but there are a few catches.
//
// The first catch is that not all immediate values are allowed. To be more
// specific, on an architecture with XLEN bits in registers (for simplicity
// let's consider XLEN=32 - 32 bit processors), it doesn't make sense to encode
// more than 31 bit immediate value to shift and though all (but one - see
// below) higher bits of immediate value are reserved to be zero.
//
// Another irregularity in immediate shift instruction encoding is the different
// in between logical and arithmetic shift. The bit differentiating logical and
// arithmetic shift is bit [30] of an instruction opcode which would correcpond
// to bit [11] of 12bit I immediate type encoding. Unfortunately this encodings
// brings a weird inconsistency when two distinct instructions identified by two
// different assembler names (srli and srai) have the same opcode but differ
// only in an immediate value bit.
//
// As different immediate value can encode different assembler instructions, we
// need them to be parsed as 2 different instructions. Consequently we are
// forced to describe this instruction opcode meta-format which is not specified
// by the architecture specification document by itself, but which allows us to
// parse bit shifts by an immediate value.
//
// The problem with differentiating logical and arithmetic shirt applies as well
// on srl and sra instructions (i.e. shift instructions accepting register
// arguments). Fortunately there we can treat the instruction opcode as 17bit
// opcode as every other bit (but bit [30]) of an immediate value is reserved to
// be zero.
func opcodeShiftImm(arithmetic bool, shiftBits uint8, mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)
	assertShiftBits(shiftBits)

	var high byte = 0
	if arithmetic {
		high = byte(1) << 6
	}

	// Shift is encoded in bits [20:(20+shiftBits)]. So we do 1<<shiftBits
	// to get 2^shiftBits. Then we subtract 1 which creates is a bit mask
	// for bits encoding values 0..(2^shiftBits)-1. We then invert the mask
	// to ensure that all other reserved bits of the actual opcode are zero.
	shiftBitMask := (uint16(1) << shiftBits) - 1
	// Then we have to shift this mask to the right place - to 20th bit of
	// opcode. As we have just high half of instruction opcode, we are
	// already shifted 16 bits. So 4 bits are remaining.
	highHalfMask := (^shiftBitMask) << 4

	return opcode.Opcode{
		Bytes: []byte{low, mid << 4, 0, high},
		Mask: []byte{
			low7Bits,
			low3Bits << 4,
			byte(highHalfMask),
			byte(highHalfMask >> 8),
		},
	}
}

// opcodeAtomic creates an opcode definition for atomic-type instruction.
//
// As atomic instructions contain acquire and release bits in bits [25] and [24]
// respectively. As those 2 bits doesn't matter for instruction opcode matching,
// we have to exclude those from opcode matching.
func opcodeAtomic(high byte, mid byte, low byte) opcode.Opcode {
	assertMask(low, low7Bits)
	assertMask(mid, low3Bits)
	assertMask(high, low5Bits)

	return opcode.Opcode{
		Bytes: []byte{low, mid << 4, 0, high << 3},
		Mask: []byte{
			low7Bits,
			low3Bits << 4,
			0,
			low5Bits << 3,
		},
	}
}

func addrAddImm(a model.Addr, imm int32) model.Addr {
	if imm >= 0 {
		return a + model.Addr(imm)
	} else {
		return a - model.Addr(-imm)
	}
}

func immConst(t immType, i instruction) expr.Const {
	imm, ok := t.parseValue(i.value)
	if !ok {
		panic(fmt.Sprintf("immediate encoding %d has no value", t))
	}
	// Immediatealways contains at most 20 bits, so 32 bits is always
	// enough.
	return expr.ConstFromInt(imm)
}

func regLoad(r reg, i instruction, w expr.Width) expr.Expr {
	num := r.regNum(i.value)
	if num == 0 {
		return expr.Zero
	}
	return expr.NewRegLoad(expr.Key(num.String()), w)
}

type binaryExprFunc func(e1, e2 expr.Expr, w expr.Width) expr.Expr

func binOpFunc(op expr.BinaryOp) binaryExprFunc {
	return func(e1, e2 expr.Expr, w expr.Width) expr.Expr {
		return expr.NewBinary(op, e1, e2, w)
	}
}

func regImmOp(f binaryExprFunc, t immType, i instruction, w expr.Width) expr.Expr {
	return f(regLoad(rs1, i, w), immConst(t, i), w)
}

func reg2Op(f binaryExprFunc, i instruction, w expr.Width) expr.Expr {
	return f(regLoad(rs1, i, w), regLoad(rs2, i, w), w)
}

func maskedRegOp(f binaryExprFunc, i instruction, bits uint8, w expr.Width) expr.Expr {
	masked := exprtools.MaskBits(regLoad(rs2, i, w), exprtools.BitCnt(bits), w)
	return f(regLoad(rs1, i, w), masked, w)
}

func regImmShift(f binaryExprFunc, i instruction, bits uint8, w expr.Width) expr.Expr {
	assertShiftBits(bits)
	mask := int32(1) << int32(bits)
	imm, _ := immTypeI.parseValue(i.value)
	immShift := expr.ConstFromInt(imm & mask)
	return f(regLoad(rs1, i, w), immShift, w)
}

func sext(e expr.Expr, signBit uint8, w expr.Width) expr.Expr {
	return exprtools.SignExtend(e, expr.ConstFromUint(signBit), w)
}

func sext32To64(e expr.Expr) expr.Expr { return sext(e, 31, expr.Width64) }

func memLoad(addr expr.Expr, w expr.Width) expr.Expr {
	return expr.NewMemLoad(MemoryKey, addr, w)
}

func memStore(e expr.Expr, addr expr.Expr, w expr.Width) expr.Effect {
	return expr.NewMemStore(e, MemoryKey, addr, w)
}

func regStore(e expr.Expr, i instruction, w expr.Width) expr.Effect {
	num := rd.regNum(i.value)
	if num == 0 {
		return nil
	}
	return expr.NewRegStore(e, expr.Key(num.String()), w)
}

func addrImmConst(t immType, i instruction, w expr.Width) expr.Const {
	imm, ok := t.parseValue(i.value)
	if !ok {
		panic(fmt.Sprintf("immediate encoding %d has no value", t))
	}
	return expr.NewConstUint(addrAddImm(i.addr, imm), w)
}

func branchCmp(
	cond expr.Condition,
	branchIfTrue bool,
	i instruction,
	w expr.Width,
) expr.Effect {
	jumpTarget := addrImmConst(immTypeB, i, w)
	nextInstr := expr.NewConstUint(i.addr+instructionLen, w)

	condTrue, condFalse := jumpTarget, nextInstr
	if !branchIfTrue {
		condTrue, condFalse = nextInstr, jumpTarget
	}

	ip := expr.NewCond(
		cond,
		regLoad(rs1, i, w),
		regLoad(rs2, i, w),
		condTrue,
		condFalse,
		w,
	)

	return expr.NewRegStore(ip, expr.IPKey, w)
}

func atomicMinMax(c expr.Condition, negate bool) binaryExprFunc {
	return func(e1, e2 expr.Expr, w expr.Width) expr.Expr {
		t, f := e1, e2
		if negate {
			t, f = e2, e1
		}

		return expr.NewCond(c, e1, e2, t, f, w)
	}
}

func atomicOp(f binaryExprFunc, i instruction, w expr.Width) []expr.Effect {
	addr := regLoad(rs1, i, w)
	ld := memLoad(addr, w)

	val := f(ld, regLoad(rs2, i, w), w)
	return []expr.Effect{
		regStore(ld, i, w),
		memStore(val, addr, w),
	}
}

func atomicOpWidth(f binaryExprFunc, i instruction, addrW, opW expr.Width) []expr.Effect {
	if opW >= width64 {
		panic(fmt.Sprintf("bug: expression too wide: %d", opW))
	}

	addr := regLoad(rs1, i, addrW)
	ld := memLoad(addr, opW)

	val := f(ld, regLoad(rs2, i, opW), opW)
	return []expr.Effect{
		regStore(sext32To64(ld), i, addrW),
		memStore(val, addr, opW),
	}
}

func csrKey(i instruction) expr.Key {
	csrNum, _ := immTypeI.parseValue(i.value)
	return expr.Key(csr(csrNum).String())
}

func csrImm(i instruction) expr.Const {
	val := uint8((i.value >> 15) & 0x1f)
	return expr.ConstFromUint(val)
}

var instructions = map[Variant]map[Extension][]*instructionType{
	Variant32: {
		extI: integer32,
		ExtM: mul32,
		ExtA: atomic32,
	},
	Variant64: {
		extI: integer64,
		ExtM: mul64,
		ExtA: atomic64,
	},
}

const (
	width8   = expr.Width8
	width16  = expr.Width16
	width32  = expr.Width32
	width64  = expr.Width64
	width128 = expr.Width128
)
