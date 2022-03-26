package expr_test

import (
	"decomp/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDynamic_Eq(t *testing.T) {
	r := require.New(t)

	d1 := expr.NewDynamic()
	d2 := expr.NewDynamic()
	r.False(d1 == d2)

	r.True(expr.Equal(d1, d1))
	r.True(expr.Equal(d2, d2))
	r.False(expr.Equal(d1, d2))
}
