package exprtools_test

import (
	"fmt"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWidthGadgetArg(t *testing.T) {
	tests := []struct {
		e   expr.Expr
		w   expr.Width
		arg expr.Expr
	}{
		{
			e: exprtools.NewWidthGadget(
				expr.NewRegLoad("reg1", expr.Width16),
				expr.Width32,
			),
			w:   expr.Width32,
			arg: expr.NewRegLoad("reg1", expr.Width16),
		},
		{
			e: exprtools.NewWidthGadget(
				expr.NewRegLoad("reg1", expr.Width64),
				expr.Width8,
			),
			w:   expr.Width8,
			arg: expr.NewRegLoad("reg1", expr.Width64),
		},
		{
			e: expr.NewRegLoad("reg1", expr.Width16),
			w: expr.Width16,
		},
		{
			e: expr.NewBinary(expr.Sub,
				expr.NewRegLoad("reg1", expr.Width64),
				expr.Zero,
				expr.Width8,
			),
			w: expr.Width8,
		},
		{
			e: expr.NewBinary(expr.Add,
				expr.NewRegLoad("reg1", expr.Width64),
				expr.NewConstUint[uint8](1, expr.Width8),
				expr.Width16,
			),
			w: expr.Width16,
		},
		{
			e: expr.NewBinary(expr.Add,
				expr.Zero,
				expr.NewRegLoad("reg1", expr.Width64),
				expr.Width32,
			),
			w: expr.Width32,
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := require.New(t)

			r.Equal(tt.w, tt.e.Width())
			arg, ok := exprtools.WidthGadgetArg(tt.e)
			if tt.arg == nil {
				r.False(ok)
			} else {
				r.True(ok)
				r.Equal(tt.arg, arg)
			}
		})
	}
}
