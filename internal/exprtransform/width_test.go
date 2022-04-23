package exprtransform_test

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetWidth(t *testing.T) {
	c1 := expr.NewConstUint[uint16](457, expr.Width16)
	c2 := expr.NewConstUint[uint32](54455554, expr.Width32)

	c3 := expr.NewConstUint[uint16](787, expr.Width16)
	c4 := expr.NewConstUint[uint32](87989030, expr.Width32)

	tests := []struct {
		name string
		e    expr.Expr
		w    expr.Width
		exp  expr.Expr
	}{
		{
			name: "identical_width",
			e:    expr.NewBinary(expr.Add, c1, c2, expr.Width16),
			w:    expr.Width16,
			exp:  expr.NewBinary(expr.Add, c1, c2, expr.Width16),
		},
		{
			name: "const_wider",
			e:    c2,
			w:    expr.Width64,
			exp:  expr.NewConstUint[uint32](54455554, expr.Width64),
		},
		{
			name: "const_narrower",
			e:    c2,
			w:    expr.Width8,
			exp:  expr.NewConstUint[uint8](54455554%256, expr.Width8),
		},
		{
			name: "binary_wider",
			e:    expr.NewBinary(expr.Add, c1, c2, expr.Width16),
			w:    expr.Width32,
			exp: exprtools.NewWidthGadget(
				expr.NewBinary(expr.Add, c1, c2, expr.Width16),
				expr.Width32,
			),
		},
		{
			name: "binary_narrower",
			e:    expr.NewBinary(expr.Add, c1, c2, expr.Width16),
			w:    expr.Width8,
			exp: exprtools.NewWidthGadget(
				expr.NewBinary(expr.Add, c1, c2, expr.Width16),
				expr.Width8,
			),
		},
		{
			name: "cond_wider",
			e:    expr.NewCond(expr.Eq, c1, c2, c3, c4, expr.Width16),
			w:    expr.Width32,
			exp: exprtools.NewWidthGadget(
				expr.NewCond(expr.Eq, c1, c2, c3, c4, expr.Width16),
				expr.Width32,
			),
		},

		{
			name: "cond_narrower",
			e:    expr.NewCond(expr.Eq, c1, c2, c3, c4, expr.Width16),
			w:    expr.Width8,
			exp: exprtools.NewWidthGadget(
				expr.NewCond(expr.Eq, c1, c2, c3, c4, expr.Width16),
				expr.Width8,
			),
		},
		{
			name: "mem_load_narrower",
			e:    expr.NewMemLoad("mem01", c1, expr.Width32),
			w:    expr.Width16,
			exp: exprtools.NewWidthGadget(
				expr.NewMemLoad("mem01", c1, expr.Width32),
				expr.Width16,
			),
		},
		{
			name: "mem_load_wider",
			e:    expr.NewMemLoad("mem01", c1, expr.Width32),
			w:    expr.Width64,
			exp: exprtools.NewWidthGadget(
				expr.NewMemLoad("mem01", c1, expr.Width32),
				expr.Width64,
			),
		},
		{
			name: "reg_load_narrower",
			e:    expr.NewRegLoad("foo1", expr.Width32),
			w:    expr.Width16,
			exp:  expr.NewRegLoad("foo1", expr.Width16),
		},
		{
			name: "reg_load_wider",
			e:    expr.NewRegLoad("foo1", expr.Width32),
			w:    expr.Width64,
			exp: exprtools.NewWidthGadget(
				expr.NewRegLoad("foo1", expr.Width32),
				expr.Width64,
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			set := exprtransform.SetWidth(tt.e, tt.w)
			require.Equal(t, tt.exp, set)
		})
	}
}

func assertSameMemory(t testing.TB, expected expr.Expr, actual expr.Expr) {
	// TODO: Find some non-deprecated way of doing this.
	require.Equal(t,
		reflect.ValueOf(&expected).Elem().InterfaceData()[1],
		reflect.ValueOf(&actual).Elem().InterfaceData()[1],
	)
}

func TestPurgeWidthGadgets(t *testing.T) {
	c1 := expr.NewConstUint[uint16](457, expr.Width16)
	c2 := expr.NewConstUint[uint32](54455554, expr.Width32)
	c3 := expr.NewConstUint[uint64](8796756564, expr.Width64)
	c4 := expr.NewConstUint[uint8](34, expr.Width8)

	tests := []struct {
		name     string
		e        expr.Expr
		exp      expr.Expr
		noChange bool
	}{
		{
			name: "no_gadgets",
			e: expr.NewBinary(expr.Sub,
				expr.NewBinary(expr.Add, c1, c2, expr.Width32),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
			noChange: true,
		},
		{
			name: "useless_output_gadget",
			e: exprtools.NewWidthGadget(expr.NewBinary(expr.Sub,
				expr.NewBinary(expr.Add, c1, c2, expr.Width32),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			), expr.Width32),
			exp: expr.NewBinary(expr.Sub,
				expr.NewBinary(expr.Add, c1, c2, expr.Width32),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
		},
		{
			name: "widening_output_gadget",
			e: exprtools.NewWidthGadget(
				expr.NewBinary(expr.Add, c1, c2, expr.Width32),
				expr.Width64,
			),
		},
		{
			name: "widening_intermediate_gadget",
			e: expr.NewBinary(expr.Sub,
				exprtools.NewWidthGadget(
					expr.NewBinary(expr.Add, c1, c2, expr.Width32),
					expr.Width64,
				),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
			exp: expr.NewBinary(expr.Sub,
				expr.NewBinary(expr.Add, c1, c2, expr.Width32),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
		},
		{
			name: "narrowing_intermediate_gadget",
			e: expr.NewBinary(expr.Sub,
				exprtools.NewWidthGadget(
					expr.NewBinary(expr.Add, c1, c2, expr.Width32),
					expr.Width16,
				),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
		},
		{
			name: "widening_intermediate_chain_big_small",
			e: expr.NewBinary(expr.Sub,
				exprtools.NewWidthGadget(exprtools.NewWidthGadget(
					expr.NewBinary(expr.Add, c1, c2, expr.Width32),
					expr.Width64), expr.Width32,
				),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
			exp: expr.NewBinary(expr.Sub,
				expr.NewBinary(expr.Add, c1, c2, expr.Width32),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
		},
		{
			name: "widening_intermediate_chain_small_big_small",
			e: expr.NewBinary(expr.Sub,
				exprtools.NewWidthGadget(exprtools.NewWidthGadget(
					exprtools.NewWidthGadget(
						expr.NewBinary(expr.Add, c1, c2, expr.Width32),
						expr.Width64,
					), expr.Width32), expr.Width64,
				),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
			exp: expr.NewBinary(expr.Sub,
				expr.NewBinary(expr.Add, c1, c2, expr.Width32),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
		},
		{
			name: "widening_intermediate_chain_smalling_chain",
			e: expr.NewBinary(expr.Sub,
				exprtools.NewWidthGadget(exprtools.NewWidthGadget(
					exprtools.NewWidthGadget(
						expr.NewBinary(expr.Add, c1, c2, expr.Width64),
						expr.Width32,
					), expr.Width16), expr.Width8,
				),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
			exp: expr.NewBinary(expr.Sub,
				exprtools.NewWidthGadget(
					expr.NewBinary(expr.Add, c1, c2, expr.Width64),
					expr.Width8,
				),
				expr.NewBinary(expr.Mod, c3, c4, expr.Width64),
				expr.Width32,
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.exp
			if expected == nil {
				expected = tt.e
			}

			purged := exprtransform.PurgeWidthGadgets(tt.e)
			require.Equal(t, expected, purged)

			if tt.noChange {
				assertSameMemory(t, expected, purged)
			}
		})
	}
}
