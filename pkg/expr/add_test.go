package expr_test

import (
	"decomp/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd_Eq(t *testing.T) {
	testEqual(t, func(e1, e2 expr.Expr) expr.Expr {
		return expr.NewAddI(expr.Width32Bit, e1, e2)
	})
	testEqual(t, func(e1, e2 expr.Expr) expr.Expr {
		return expr.NewAddI(expr.Width64Bit, e1, e2)
	})

	t.Run("not_equal_bit_width", func(t *testing.T) {
		c1 := expr.NewConst([]byte{5})
		c2 := expr.NewConst([]byte{7})

		a1 := expr.NewAddI(expr.Width32Bit, c1, c2)
		a2 := expr.NewAddI(expr.Width64Bit, c1, c2)

		require.True(t, expr.Equal(a1, a2))
	})
}
