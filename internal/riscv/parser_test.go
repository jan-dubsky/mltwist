package riscv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func pow2(n int) int {
	p := 1
	for i := 0; i < n; i++ {
		p *= 2
	}
	return p
}

func TestNewParser(t *testing.T) {
	allExts := make([]Extension, 0, extEnd-1)
	for e := extI + 1; e < extEnd; e++ {
		allExts = append(allExts, e)
	}

	// generate powerset of exts
	exts := make([][]Extension, pow2(len(allExts)))
	for i := 0; i < len(exts); i++ {
		for j := 0; j < len(allExts); j++ {
			if i&pow2(j) != 0 {
				exts[i] = append(exts[i], allExts[j])
			}
		}
	}

	for v := Variant32; v < variantEnd; v++ {
		v := v
		t.Run(fmt.Sprintf("variant_%d", v), func(t *testing.T) {
			t.Run("integer_instructions_invalid", func(t *testing.T) {
				require.Panics(t, func() {
					NewParser(v, extI)
				})
			})
			t.Run("no_extension_valid", func(_ *testing.T) {
				NewParser(v)
			})

			for _, es := range exts {
				t.Run(fmt.Sprintf("ext_%v", es), func(_ *testing.T) {
					NewParser(v)
				})
			}
		})
	}
}
