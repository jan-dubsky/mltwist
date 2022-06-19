package exprtransform_test

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	c1 := expr.ConstFromUint[uint32](55)
	c2 := expr.ConstFromUint[uint8](22)

	tests := []struct {
		name string
		e1   expr.Expr
		e2   expr.Expr
		eq   bool
	}{{
		name: "const_eq",
		e1:   c1,
		e2:   c1,
		eq:   true,
	}, {
		name: "const_neq_width",
		e1:   c1,
		e2:   expr.ConstFromUint[uint16](55),
		eq:   false,
	}, {
		name: "const_neq_value",
		e1:   expr.ConstFromUint[uint32](2139032992),
		e2:   expr.ConstFromUint[uint32](2139032992 + 1),
		eq:   false,
	}, {
		name: "mem_load_eq",
		e1:   expr.NewMemLoad("mem", c1, expr.Width32),
		e2:   expr.NewMemLoad("mem", c1, expr.Width32),
		eq:   true,
	}, {
		name: "mem_load_neq_mem",
		e1:   expr.NewMemLoad("mem", c1, expr.Width32),
		e2:   expr.NewMemLoad("mem2", c1, expr.Width32),
		eq:   false,
	}, {
		name: "mem_load_neq_width",
		e1:   expr.NewMemLoad("mem", c1, expr.Width16),
		e2:   expr.NewMemLoad("mem", c1, expr.Width32),
		eq:   false,
	}, {
		name: "reg_load_eq",
		e1:   expr.NewRegLoad("foo1", expr.Width32),
		e2:   expr.NewRegLoad("foo1", expr.Width32),
		eq:   true,
	}, {
		name: "reg_load_neq_width",
		e1:   expr.NewRegLoad("foo1", expr.Width32),
		e2:   expr.NewRegLoad("foo1", expr.Width16),
		eq:   false,
	}, {
		name: "reg_load_neq_reg",
		e1:   expr.NewRegLoad("foo1", expr.Width32),
		e2:   expr.NewRegLoad("foo2", expr.Width32),
		eq:   false,
	}, {
		name: "binary_eq",
		e1:   expr.NewBinary(expr.Add, c1, c2, expr.Width32),
		e2:   expr.NewBinary(expr.Add, c1, c2, expr.Width32),
		eq:   true,
	}, {
		name: "binary_neq_width",
		e1:   expr.NewBinary(expr.Add, c1, c2, expr.Width32),
		e2:   expr.NewBinary(expr.Add, c1, c2, expr.Width16),
		eq:   false,
	}, {
		name: "binary_neq_arg1",
		e1: expr.NewBinary(expr.Add,
			expr.NewConstUint[uint8](54, expr.Width128),
			c2,
			expr.Width32,
		),
		e2: expr.NewBinary(expr.Add, c1, c2, expr.Width16),
		eq: false,
	}, {
		name: "binary_neq_arg2",
		e1: expr.NewBinary(expr.Add,
			c1,
			expr.ConstFromUint[uint8](22),
			expr.Width32,
		),
		e2: expr.NewBinary(expr.Add,
			c1,
			expr.ConstFromUint[uint8](23),
			expr.Width16,
		),
		eq: false,
	}, {
		name: "cond_eq",
		e1:   expr.NewLess(c1, c2, c1, c2, expr.Width16),
		e2:   expr.NewLess(c1, c2, c1, c2, expr.Width16),
		eq:   true,
	}, {
		name: "cond_neq_width",
		e1:   expr.NewLess(c1, c2, c1, c2, expr.Width16),
		e2:   expr.NewLess(c1, c2, c1, c2, expr.Width32),
		eq:   false,
	}, {
		name: "cond_neq_arg1",
		e1:   expr.NewLess(c1, c2, c1, c2, expr.Width16),
		e2:   expr.NewLess(c2, c2, c1, c2, expr.Width16),
		eq:   false,
	}, {
		name: "cond_neq_arg2",
		e1:   expr.NewLess(c1, c2, c1, c2, expr.Width16),
		e2:   expr.NewLess(c1, c1, c1, c2, expr.Width16),
		eq:   false,
	}, {
		name: "cond_neq_true",
		e1:   expr.NewLess(c1, c2, c1, c2, expr.Width16),
		e2:   expr.NewLess(c1, c2, c2, c2, expr.Width16),
		eq:   false,
	}, {
		name: "cond_neq_false",
		e1:   expr.NewLess(c1, c2, c1, c2, expr.Width16),
		e2:   expr.NewLess(c1, c2, c1, c1, expr.Width16),
		eq:   false,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ok := exprtransform.Equal(tt.e1, tt.e2)
			require.Equal(t, tt.eq, ok)
		})
	}
}
