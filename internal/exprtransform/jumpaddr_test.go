package exprtransform_test

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJumpAddr(t *testing.T) {
	tests := []struct {
		name  string
		e     expr.Expr
		addrs []expr.Expr
	}{{
		name: "const_offset_branch",
		e: expr.NewCond(
			expr.Eq,
			expr.NewRegLoad("r1", expr.Width32),
			expr.Zero,
			expr.NewConstUint[uint32](0x1000, expr.Width32),
			expr.NewConstUint[uint32](0x1000+52, expr.Width32),
			expr.Width32,
		),
		addrs: []expr.Expr{
			expr.NewConstUint[uint32](0x1000, expr.Width32),
			expr.NewConstUint[uint32](0x1000+52, expr.Width32),
		},
	}, {
		name: "const_offset_branch_with_width",
		e: expr.NewCond(
			expr.Eq,
			expr.NewRegLoad("r1", expr.Width8),
			expr.Zero,
			expr.NewConstUint[uint32](0x1000, expr.Width64),
			expr.NewConstUint[uint32](0x1000+52, expr.Width16),
			expr.Width32,
		),
		addrs: []expr.Expr{
			expr.NewConstUint[uint32](0x1000, expr.Width32),
			expr.NewConstUint[uint32](0x1000+52, expr.Width32),
		},
	}, {
		name: "cond_followed_by_calculation",
		e: expr.NewBinary(expr.Add,
			expr.NewCond(
				expr.Eq,
				expr.NewRegLoad("r1", expr.Width32),
				expr.Zero,
				expr.NewConstUint[uint32](0x1000, expr.Width32),
				expr.NewConstUint[uint32](0x1000+52, expr.Width32),
				expr.Width32,
			),
			expr.NewConstUint[uint32](98, expr.Width8),
			expr.Width64,
		),
		addrs: []expr.Expr{
			expr.NewConstUint[uint32](0x1000+98, expr.Width64),
			expr.NewConstUint[uint32](0x1000+52+98, expr.Width64),
		},
	}, {
		name: "cross_product",
		e: expr.NewBinary(expr.Add,
			expr.NewCond(
				expr.Eq,
				expr.NewRegLoad("r1", expr.Width32),
				expr.Zero,
				expr.NewConstUint[uint32](0x1000, expr.Width32),
				expr.NewConstUint[uint32](0x10aa, expr.Width32),
				expr.Width32,
			),
			expr.NewCond(
				expr.Ltu,
				expr.NewRegLoad("r2", expr.Width32),
				expr.NewRegLoad("r3", expr.Width32),
				expr.NewConstUint[uint32](0xff, expr.Width32),
				expr.NewConstUint[uint32](0x888, expr.Width32),
				expr.Width32,
			),
			expr.Width64,
		),
		addrs: []expr.Expr{
			expr.NewConstUint[uint32](0x1000+0xff, expr.Width64),
			expr.NewConstUint[uint32](0x1000+0x888, expr.Width64),
			expr.NewConstUint[uint32](0x10aa+0xff, expr.Width64),
			expr.NewConstUint[uint32](0x10aa+0x888, expr.Width64),
		},
	}, {
		name: "mem_ref_cross",
		e: expr.NewBinary(expr.Add,
			expr.NewCond(
				expr.Eq,
				expr.NewRegLoad("r1", expr.Width32),
				expr.Zero,
				expr.NewConstUint[uint32](0x1000, expr.Width32),
				expr.NewConstUint[uint32](0x10aa, expr.Width32),
				expr.Width32,
			),
			expr.NewCond(
				expr.Ltu,
				expr.NewRegLoad("r2", expr.Width32),
				expr.NewRegLoad("r3", expr.Width32),
				expr.NewConstUint[uint32](0xff, expr.Width32),
				expr.NewMemLoad(
					"mem",
					expr.NewConstUint[uint32](0x8765, expr.Width32),
					expr.Width16,
				),
				expr.Width32,
			),
			expr.Width64,
		),
		addrs: []expr.Expr{
			expr.NewConstUint[uint32](0x1000+0xff, expr.Width64),
			expr.NewBinary(expr.Add,
				expr.NewConstUint[uint32](0x1000, expr.Width32),
				expr.NewMemLoad(
					"mem",
					expr.NewConstUint[uint32](0x8765, expr.Width32),
					expr.Width16,
				),
				expr.Width64,
			),
			expr.NewConstUint[uint32](0x10aa+0xff, expr.Width64),
			expr.NewBinary(expr.Add,
				expr.NewConstUint[uint32](0x10aa, expr.Width32),
				expr.NewMemLoad(
					"mem",
					expr.NewConstUint[uint32](0x8765, expr.Width32),
					expr.Width16,
				),
				expr.Width64,
			),
		},
	},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			addrs := exprtransform.JumpAddrs(tt.e)
			require.Equal(t, tt.addrs, addrs)
		})
	}
}
