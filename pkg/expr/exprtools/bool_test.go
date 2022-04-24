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
			e:   expr.ConstFromUint[uint8](8),
			exp: expr.ConstFromUint[uint8](1),
		},
		{
			e:   expr.ConstFromUint[uint8](0),
			exp: expr.ConstFromUint[uint8](0),
		},
		{
			e:   expr.ConstFromUint[uint16](0xba00),
			exp: expr.ConstFromUint[uint8](1),
		},
		{
			e:   expr.ConstFromUint[uint32](0),
			exp: expr.ConstFromUint[uint8](0),
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
			e:   expr.ConstFromUint[uint8](8),
			exp: expr.ConstFromUint[uint8](0),
		},
		{
			e:   expr.ConstFromUint[uint8](0),
			exp: expr.ConstFromUint[uint8](1),
		},
		{
			e:   expr.ConstFromUint[uint64](0xfba462839bacaef1),
			exp: expr.ConstFromUint[uint8](0),
		},
		{
			e:   expr.ConstFromUint[uint16](0x1700),
			exp: expr.ConstFromUint[uint8](0),
		},
		{
			e:   expr.ConstFromUint[uint32](0),
			exp: expr.ConstFromUint[uint8](1),
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
			e:   expr.ConstFromUint[uint8](8),
			t:   expr.ConstFromUint[uint32](0xfba64bc5),
			f:   expr.ConstFromUint[uint8](1),
			w:   expr.Width32,
			exp: expr.ConstFromUint[uint32](0xfba64bc5),
		},
		{
			e:   expr.ConstFromUint[uint8](0),
			t:   expr.ConstFromUint[uint32](0xfba64bc5),
			f:   expr.ConstFromUint[uint8](1),
			w:   expr.Width16,
			exp: expr.ConstFromUint[uint16](1),
		},
		{
			e:   expr.ConstFromUint[uint8](8),
			t:   expr.ConstFromUint[uint32](0xfba64bc5),
			f:   expr.ConstFromUint[uint8](1),
			w:   expr.Width16,
			exp: expr.ConstFromUint[uint16](0x4bc5),
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
