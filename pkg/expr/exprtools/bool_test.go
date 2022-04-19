package exprtools_test

import (
	"fmt"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	tests := []struct {
		e   expr.Expr
		exp expr.Expr
	}{
		{
			e:   expr.NewConstUint[uint8](8, expr.Width8),
			exp: expr.NewConstUint[uint8](1, expr.Width8),
		},
		{
			e:   expr.NewConstUint[uint8](0, expr.Width8),
			exp: expr.NewConstUint[uint8](0, expr.Width8),
		},
		{
			e:   expr.NewConstUint[uint16](0xba00, expr.Width16),
			exp: expr.NewConstUint[uint8](1, expr.Width8),
		},
		{
			e:   expr.NewConstUint[uint8](0, expr.Width32),
			exp: expr.NewConstUint[uint8](0, expr.Width8),
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			e := exprtools.Bool(tt.e)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestNot(t *testing.T) {
	tests := []struct {
		e   expr.Expr
		exp expr.Expr
	}{
		{
			e:   expr.NewConstUint[uint8](8, expr.Width8),
			exp: expr.NewConstUint[uint8](0, expr.Width8),
		},
		{
			e:   expr.NewConstUint[uint8](0, expr.Width8),
			exp: expr.NewConstUint[uint8](1, expr.Width8),
		},
		{
			e:   expr.NewConstUint[uint64](0xfba462839bacaef1, expr.Width64),
			exp: expr.NewConstUint[uint8](0, expr.Width8),
		},
		{
			e:   expr.NewConstUint[uint16](0x1700, expr.Width16),
			exp: expr.NewConstUint[uint8](0, expr.Width8),
		},
		{
			e:   expr.NewConstUint[uint8](0, expr.Width32),
			exp: expr.NewConstUint[uint8](1, expr.Width8),
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			e := exprtools.Not(tt.e)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}

func TestBoolCond(t *testing.T) {
	tests := []struct {
		e   expr.Expr
		t   expr.Expr
		f   expr.Expr
		w   expr.Width
		exp expr.Expr
	}{
		{
			e:   expr.NewConstUint[uint8](8, expr.Width8),
			t:   expr.NewConstUint[uint32](0xfba64bc5, expr.Width32),
			f:   expr.NewConstUint[uint8](1, expr.Width8),
			w:   expr.Width32,
			exp: expr.NewConstUint[uint32](0xfba64bc5, expr.Width32),
		},
		{
			e:   expr.NewConstUint[uint8](0, expr.Width8),
			t:   expr.NewConstUint[uint32](0xfba64bc5, expr.Width32),
			f:   expr.NewConstUint[uint8](1, expr.Width8),
			w:   expr.Width16,
			exp: expr.NewConstUint[uint16](1, expr.Width16),
		},
		{
			e:   expr.NewConstUint[uint8](8, expr.Width8),
			t:   expr.NewConstUint[uint32](0xfba64bc5, expr.Width32),
			f:   expr.NewConstUint[uint8](1, expr.Width8),
			w:   expr.Width16,
			exp: expr.NewConstUint[uint32](0x4bc5, expr.Width16),
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			e := exprtools.BoolCond(tt.e, tt.t, tt.f, tt.w)
			require.Equal(t, tt.exp, exprtransform.ConstFold(e))
		})
	}
}
