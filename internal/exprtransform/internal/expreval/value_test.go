package expreval_test

import (
	"encoding/binary"
	"mltwist/internal/exprtransform/internal/expreval"
	"mltwist/pkg/expr"
	"testing"

	"github.com/stretchr/testify/require"
)

func valUint(v uint64, w expr.Width) expreval.Value {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, v)
	return expreval.NewValue(bs).SetWidth(w)
}

func TestValue_SetWidth(t *testing.T) {
	tests := []struct {
		name string
		v    expreval.Value
		w    expr.Width
		exp  expreval.Value
	}{
		{
			name: "same_size",
			v:    expreval.NewValue([]byte{1, 2, 3, 4}),
			w:    expr.Width32,
			exp:  expreval.NewValue([]byte{1, 2, 3, 4}),
		},
		{
			name: "cut_size",
			v:    expreval.NewValue([]byte{1, 2, 3, 4}),
			w:    expr.Width16,
			exp:  expreval.NewValue([]byte{1, 2}),
		},
		{
			name: "zero_extend",
			v:    expreval.NewValue([]byte{1, 2, 3, 4}),
			w:    expr.Width64,
			exp:  expreval.NewValue([]byte{1, 2, 3, 4, 0, 0, 0, 0}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tt.exp, tt.v.SetWidth(tt.w))
		})
	}
}
