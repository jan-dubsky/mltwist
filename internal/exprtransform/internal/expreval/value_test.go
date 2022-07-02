package expreval

import (
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValue_setWidth(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		w    expr.Width
		exp  Value
	}{{
		name: "same_size",
		v:    newValue([]byte{1, 2, 3, 4}),
		w:    expr.Width32,
		exp:  newValue([]byte{1, 2, 3, 4}),
	}, {
		name: "cut_size",
		v:    newValue([]byte{1, 2, 3, 4}),
		w:    expr.Width16,
		exp:  newValue([]byte{1, 2}),
	}, {
		name: "zero_extend",
		v:    newValue([]byte{1, 2, 3, 4}),
		w:    expr.Width64,
		exp:  newValue([]byte{1, 2, 3, 4, 0, 0, 0, 0}),
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tt.exp, tt.v.setWidth(tt.w))
		})
	}
}
