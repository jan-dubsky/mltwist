package state

import (
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestState_Apply(t *testing.T) {
	tests := []struct {
		name       string
		ef         expr.Effect
		notApplied bool
		regChanged bool
		memChanged bool
	}{{
		name: "reg_store",
		ef: expr.NewRegStore(
			expr.ConstFromUint[uint32](0x5678abcd),
			"r1",
			expr.Width64,
		),
		regChanged: true,
	}, {
		name: "mem_store",
		ef: expr.NewMemStore(
			expr.ConstFromUint[uint32](0x5678abcd),
			"r1",
			expr.Zero,
			expr.Width64,
		),
		memChanged: true,
	}, {
		name: "mem_store_unknown_address",
		ef: expr.NewMemStore(
			expr.ConstFromUint[uint32](0x5678abcd),
			"r1",
			expr.NewRegLoad("r1", expr.Width16),
			expr.Width64,
		),
		notApplied: true,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			s := New()

			ok := s.Apply(tt.ef)
			if tt.notApplied {
				r.False(ok)
				r.Equal(New(), s)
				return
			}

			r.True(ok)
			if tt.regChanged {
				r.NotEmpty(s.Regs)
			} else {
				r.Empty(s.Regs.Values())
			}

			if tt.memChanged {
				r.NotEmpty(s.Mems)
			} else {
				r.Empty(s.Mems)
			}
		})
	}
}
