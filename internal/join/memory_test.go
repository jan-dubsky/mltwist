package join

import (
	"fmt"
	"mltwist/pkg/expr"
	"mltwist/pkg/expr/exprtools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCutExpr_Cuts(t *testing.T) {
	tests := []struct {
		e     cutExpr
		l     expr.Width
		begin expr.Width
		end   expr.Width
		panic bool
	}{
		{
			e:     cutExpr{nil, 0, 8},
			l:     5,
			begin: 3,
			end:   5,
		},
		{
			e:     cutExpr{nil, 0, 8},
			l:     8,
			begin: 0,
			end:   8,
		},
		{
			e:     cutExpr{nil, 0, 8},
			l:     0,
			begin: 8,
			end:   0,
		},
		{
			e:     cutExpr{nil, 0, 8},
			l:     9,
			panic: true,
		},
		{
			e:     cutExpr{nil, 7, 13},
			l:     4,
			begin: 9,
			end:   11,
		},
		{
			e:     cutExpr{nil, 7, 13},
			l:     5,
			begin: 8,
			end:   12,
		},
		{
			e:     cutExpr{nil, 7, 13},
			l:     6,
			begin: 7,
			end:   13,
		},
		{
			e:     cutExpr{nil, 7, 13},
			l:     7,
			panic: true,
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			r := require.New(t)
			orig := tt.e

			if tt.panic {
				r.Panics(func() {
					_ = tt.e.cutBegin(tt.l)
				})
				r.Equal(orig, tt.e)

				r.Panics(func() {
					_ = tt.e.cutEnd(tt.l)
				})
				r.Equal(orig, tt.e)

				return
			}

			ending := tt.e.cutBegin(tt.l)
			r.Equal(tt.e.end, ending.end)
			r.Equal(tt.begin, ending.begin)
			r.Equal(orig, tt.e)

			beginning := tt.e.cutEnd(tt.l)
			r.Equal(tt.e.begin, beginning.begin)
			r.Equal(tt.end, beginning.end)
			r.Equal(orig, tt.e)
		})
	}
}

func TestCutExpr_Expr(t *testing.T) {
	tests := []struct {
		e   cutExpr
		exp expr.Expr
	}{
		{
			e: cutExpr{
				ex:    expr.ConstFromUint[uint8](56),
				begin: 0,
				end:   4,
			},
			exp: expr.ConstFromUint[uint32](56),
		},
		{
			e: cutExpr{
				ex:    expr.ConstFromUint[uint32](0x12345678),
				begin: 0,
				end:   2,
			},
			exp: expr.ConstFromUint[uint16](0x5678),
		},
		{
			e: cutExpr{
				ex:    expr.ConstFromUint[uint32](0x12345678),
				begin: 2,
				end:   4,
			},
			exp: exprtools.NewWidthGadget(expr.NewBinary(expr.Rsh,
				expr.ConstFromUint[uint32](0x12345678),
				expr.ConstFromUint[uint16](16),
				expr.Width32,
			), expr.Width16),
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			require.Equal(t, tt.exp, tt.e.expr())
		})
	}
}
