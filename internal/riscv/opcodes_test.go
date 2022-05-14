package riscv

import (
	"fmt"
	"mltwist/internal/opcode"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type hashableOpcode struct {
	Bytes string
	Mask  string
}

func newHashableOpcode(o opcode.Opcode) hashableOpcode {
	return hashableOpcode{
		Bytes: fmt.Sprintf("%x", o.Bytes),
		Mask:  fmt.Sprintf("%x", o.Mask),
	}
}

func assertValidOpcode(t testing.TB, xlenBytes uint8, instrs []*instructionOpcode) {
	for _, ins := range instrs {
		require.NoError(t, ins.validate(xlenBytes), ins.Name())
	}
}

func assertUniqueOpcode(t testing.TB, instrs []*instructionOpcode) {
	opcodeSet := make(map[hashableOpcode]struct{}, len(instrs))
	for i, ins := range instrs {
		o := ins.Opcode()
		h := newHashableOpcode(o)
		if _, ok := opcodeSet[h]; ok {
			require.Failf(t, "opcode is not unique", "%d: %#v (%#v)",
				i, o, ins)
		}
		opcodeSet[h] = struct{}{}
	}
}

func assertUniqueNames(t testing.TB, instrs []*instructionOpcode) {
	m := make(map[string]*instructionOpcode, len(instrs))
	for _, ins := range instrs {
		_, ok := m[ins.name]
		require.False(t, ok)
		m[ins.name] = ins
	}
}

func TestInstructions(t *testing.T) {
	for arch, exts := range instructions {
		exts := exts
		t.Run(fmt.Sprintf("arch_%v", arch), func(t *testing.T) {
			allInstrs := make([]*instructionOpcode, 0)
			for _, instrs := range exts {
				allInstrs = append(allInstrs, instrs...)
			}

			assertUniqueOpcode(t, allInstrs)
			assertUniqueNames(t, allInstrs)

			for ext, instrs := range exts {
				instrs := instrs
				t.Run(fmt.Sprintf("ext_%v", ext), func(t *testing.T) {
					xlenBytes := 4
					if arch == Variant64 {
						xlenBytes = 8
					}

					assertValidOpcode(t, uint8(xlenBytes), instrs)
				})
			}
		})
	}
}

func opcodeMap(t testing.TB, opcs []*instructionOpcode) map[string]*instructionOpcode {
	m := make(map[string]*instructionOpcode)
	for _, o := range opcs {
		_, ok := m[o.name]
		require.False(t, ok)
		m[o.name] = o
	}

	return m
}

func TestRV32ToRV64(t *testing.T) {
	type instrSet struct {
		name string
		rv32 []*instructionOpcode
		rv64 []*instructionOpcode
	}

	lists := []instrSet{
		{
			name: "integer",
			rv32: integer32,
			rv64: integer64,
		}, {
			name: "mul",
			rv32: mul32,
			rv64: mul64,
		},
	}

	for _, set := range lists {
		set := set
		t.Run(set.name, func(t *testing.T) {
			r := require.New(t)
			rv32 := opcodeMap(t, set.rv32)
			rv64 := opcodeMap(t, set.rv64)

			for name, opc32 := range rv32 {
				opc64, ok := rv64[name]
				r.True(ok, "RV64 version of %q not found.", name)

				r.Equal(opc32.inputRegCnt, opc64.inputRegCnt)
				r.Equal(opc32.hasOutputReg, opc64.hasOutputReg)
				r.Equal(opc32.immediate, opc64.immediate)
				r.Equal(opc32.instrType, opc64.instrType)
			}

			for name := range rv64 {
				if _, ok := rv32[name]; ok {
					continue
				}

				// RV64 relates loads and stores.
				if name == "sd" || name == "ld" || name == "lwu" {
					continue
				}

				r.True(strings.HasSuffix(name, "w"),
					"RV64-only instruction does not have w suffix: %s",
					name,
				)
			}
		})
	}
}
