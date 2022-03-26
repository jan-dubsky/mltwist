package expr_test

import (
	"decomp/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

type expr2Func func(e1 expr.Expr, e2 expr.Expr) expr.Expr

func testEqual(t testing.TB, f expr2Func) {
	// TODO: Drop this.
	return

	r := require.New(t)

	c := expr.Zero
	d := expr.NewDynamic()

	e1 := f(c, d)
	e2 := f(c, d)

	r.True(expr.Equal(e1, e2))
}
