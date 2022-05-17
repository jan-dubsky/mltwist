package exprtransform_test

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstFold(t *testing.T) {
	c1 := expr.ConstFromUint[uint8](25)
	c2 := expr.ConstFromInt[int16](-1)

	tests := []struct {
		name string
		e    expr.Expr
		exp  expr.Expr
	}{{
		name: "add_const_no_width_change",
		e: expr.NewBinary(expr.Add,
			expr.ConstFromUint[uint32](23),
			expr.ConstFromUint[uint32](0xffff),
			expr.Width32,
		),
		exp: expr.ConstFromUint[uint32](0xffff + 23),
	}, {
		name: "add_const_width",
		e: expr.NewBinary(expr.Add,
			expr.ConstFromUint[uint8](23),
			expr.ConstFromUint[uint16](0xffff),
			expr.Width32,
		),
		exp: expr.ConstFromUint[uint32](0xffff + 23),
	}, {
		name: "add_const_multilayer",
		e: expr.NewBinary(expr.Add,
			expr.NewBinary(expr.Add, c1, c2, expr.Width32),
			expr.ConstFromUint[uint16](56899),
			expr.Width32,
		),
		exp: expr.ConstFromUint[uint32](25 + 0xffff + 56899),
	}, {
		name: "simplify_sub",
		e: expr.NewBinary(expr.Add,
			expr.NewBinary(expr.Sub, c1, c2, expr.Width32),
			expr.NewRegLoad("foo1", expr.Width16),
			expr.Width32,
		),
		exp: expr.NewBinary(expr.Add,
			expr.ConstFromUint[uint32](26+0xffff0000),
			expr.NewRegLoad("foo1", expr.Width16),
			expr.Width32,
		),
	}, {
		name: "eval_condition__rue_no_cond_width",
		e: expr.NewCond(
			expr.Eq,
			expr.NewBinary(expr.Lsh,
				expr.ConstFromUint[uint8](5),
				expr.NewBinary(expr.Sub, c1, c2, expr.Width8),
				expr.Width32,
			),
			expr.ConstFromUint[uint32](5<<26),
			expr.ConstFromUint[uint16](42),
			expr.Zero,
			expr.Width32,
		),
		exp: expr.ConstFromUint[uint32](42),
	}, {
		name: "eval_condition_false_no_cond_width",
		e: expr.NewCond(
			expr.Eq,
			expr.NewBinary(expr.Lsh,
				expr.ConstFromUint[uint8](5),
				expr.NewBinary(expr.Sub, c1, c2, expr.Width8),
				expr.Width32,
			),
			expr.ConstFromUint[uint32](5<<26),
			expr.ConstFromUint[uint16](42),
			expr.Zero,
			expr.Width32,
		),
		exp: expr.ConstFromUint[uint32](42),
	}, {
		name: "eval_condition_cond_width",
		e: expr.NewCond(expr.Eq,
			expr.NewBinary(expr.Lsh,
				expr.ConstFromUint[uint8](5),
				expr.NewBinary(expr.Sub, c1, c2, expr.Width8),
				expr.Width32,
			),
			expr.Zero,
			expr.ConstFromUint[uint8](42),
			expr.Zero,
			expr.Width16,
		),
		exp: expr.ConstFromUint[uint16](42),
	}, {
		name: "simplify_condition",
		e: expr.NewCond(expr.Lts,
			expr.NewBinary(expr.Mul,
				expr.ConstFromUint[uint8](5),
				expr.NewBinary(expr.Or,
					expr.ConstFromUint[uint16](0x1234),
					expr.ConstFromUint[uint32](0xdbca8765),
					expr.Width64,
				),
				expr.Width32,
			),
			expr.NewBinary(expr.Div,
				expr.NewBinary(expr.Mod,
					expr.ConstFromUint[uint32](146),
					expr.ConstFromInt[int16](13),
					expr.Width8,
				),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.ConstFromUint[uint16](666),
			expr.ConstFromInt[int32](-1),
			expr.Width16,
		),
		exp: expr.NewCond(expr.Lts,
			expr.ConstFromUint(uint32(((0x1234|0xdbca8765)*5)&0xffffffff)),
			expr.NewBinary(expr.Div,
				expr.ConstFromUint[uint8](146%13),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.ConstFromUint[uint16](666),
			expr.ConstFromInt[int32](-1),
			expr.Width16,
		),
	}, {
		name: "condition_simplify_args",
		e: expr.NewCond(expr.Ltu,
			expr.ConstFromUint[uint16](324),
			expr.NewMemLoad("mem01",
				expr.ConstFromUint[uint64](0xf3920bada83),
				expr.Width8,
			),
			expr.NewBinary(expr.Div,
				expr.NewBinary(expr.Mod,
					expr.ConstFromUint[uint32](146),
					expr.ConstFromInt[int16](13),
					expr.Width8,
				),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.NewBinary(expr.Mul,
				expr.ConstFromUint[uint8](5),
				expr.NewBinary(expr.Or,
					expr.ConstFromUint[uint16](0x1234),
					expr.ConstFromUint[uint32](0xdbca8765),
					expr.Width64,
				),
				expr.Width32,
			),
			expr.Width64,
		),
		exp: expr.NewCond(expr.Ltu,
			expr.ConstFromUint[uint16](324),
			expr.NewMemLoad("mem01",
				expr.ConstFromUint[uint64](0xf3920bada83),
				expr.Width8,
			),
			expr.NewBinary(expr.Div,
				expr.ConstFromUint[uint8](146%13),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.ConstFromUint(uint32(((0x1234|0xdbca8765)*5)&0xffffffff)),
			expr.Width64,
		),
	}, {
		name: "eval_condition_output_width_change",
		e: expr.NewCond(expr.Lts,
			expr.ConstFromUint[uint16](323),
			expr.ConstFromUint[uint16](324),
			expr.NewBinary(expr.Div,
				expr.NewBinary(expr.Mod,
					expr.ConstFromUint[uint32](146),
					expr.ConstFromInt[int16](13),
					expr.Width8,
				),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.NewBinary(expr.Mul,
				expr.ConstFromUint[uint8](5),
				expr.NewBinary(expr.Or,
					expr.ConstFromUint[uint16](0x1234),
					expr.ConstFromUint[uint32](0xdbca8765),
					expr.Width64,
				),
				expr.Width32,
			),
			expr.Width64,
		),
		exp: exprtools.NewWidthGadget(
			expr.NewBinary(expr.Div,
				expr.ConstFromUint[uint8](146%13),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.Width64,
		),
	}, {
		name: "eval_condition_intermediate_width_change_necessary",
		e: expr.NewBinary(expr.Add,
			expr.NewCond(expr.Lts,
				expr.ConstFromUint[uint16](323),
				expr.ConstFromUint[uint16](324),
				expr.NewBinary(expr.Div,
					expr.ConstFromUint[uint16](98),
					expr.NewRegLoad("foo2", expr.Width16),
					expr.Width64,
				),
				expr.ConstFromUint[uint8](45),
				expr.Width32,
			),
			expr.ConstFromUint[uint8](64),
			expr.Width64,
		),
		exp: expr.NewBinary(expr.Add,
			exprtools.NewWidthGadget(
				expr.NewBinary(expr.Div,
					expr.ConstFromUint[uint16](98),
					expr.NewRegLoad("foo2", expr.Width16),
					expr.Width64,
				),
				expr.Width32,
			),
			expr.ConstFromUint[uint8](64),
			expr.Width64,
		),
	}, {
		name: "eval_condition_intermediate_width_change_non_necessary",
		e: expr.NewBinary(expr.Add,
			expr.NewCond(expr.Ltu,
				expr.ConstFromUint[uint16](323),
				expr.ConstFromUint[uint16](324),
				expr.NewBinary(expr.Div,
					expr.ConstFromUint[uint16](98),
					expr.NewRegLoad("foo2", expr.Width16),
					expr.Width16,
				),
				expr.ConstFromUint[uint8](45),
				expr.Width32,
			),
			expr.ConstFromUint[uint8](64),
			expr.Width64,
		),
		exp: expr.NewBinary(expr.Add,
			expr.NewBinary(expr.Div,
				expr.ConstFromUint[uint16](98),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width16,
			),
			expr.ConstFromUint[uint8](64),
			expr.Width64,
		),
	}, {
		name: "memory_address",
		e: expr.NewMemLoad("foo1",
			expr.NewBinary(expr.Div,
				expr.NewBinary(expr.Mod,
					expr.ConstFromUint[uint32](146),
					expr.ConstFromInt[int16](13),
					expr.Width8,
				),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.Width64,
		),
		exp: expr.NewMemLoad("foo1",
			expr.NewBinary(expr.Div,
				expr.ConstFromUint[uint8](146%13),
				expr.NewRegLoad("foo2", expr.Width16),
				expr.Width32,
			),
			expr.Width64,
		),
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := exprtransform.ConstFold(tt.e)
			require.Equal(t, tt.exp, res)
		})
	}
}
