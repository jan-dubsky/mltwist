package expreval

import (
	"decomp/pkg/expr"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func valUint(v uint64, w expr.Width) value {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, v)
	return value(bs).setWidth(w)
}

func TestValue_SetWidth(t *testing.T) {
	tests := []struct {
		name string
		v    value
		w    expr.Width
		exp  value
	}{
		{
			name: "same_size",
			v:    value{1, 2, 3, 4},
			w:    expr.Width32,
			exp:  value{1, 2, 3, 4},
		},
		{
			name: "cut_size",
			v:    value{1, 2, 3, 4},
			w:    expr.Width16,
			exp:  value{1, 2},
		},
		{
			name: "zero_extend",
			v:    value{1, 2, 3, 4},
			w:    expr.Width64,
			exp:  value{1, 2, 3, 4, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tt.exp, tt.v.setWidth(tt.w))
		})
	}
}
