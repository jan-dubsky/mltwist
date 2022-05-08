package interval_test

import (
	"mltwist/internal/state/interval"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterval(t *testing.T) {
	tests := []struct {
		begin    int
		end      int
		len      int
		overlaps map[int]bool
	}{
		{
			begin: 1,
			end:   5,
			len:   4,
			overlaps: map[int]bool{
				0:  false,
				1:  true,
				2:  true,
				3:  true,
				4:  true,
				5:  false,
				90: false,
			},
		}, {
			begin: 5,
			end:   5,
			len:   0,
			overlaps: map[int]bool{
				4: false,
				5: false,
				6: false,
			},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := require.New(t)

			intv := interval.New(tt.begin, tt.end)
			r.Equal(tt.begin, intv.Begin())
			r.Equal(tt.end, intv.End())
			r.Equal(tt.len, intv.Len())

			for val, exp := range tt.overlaps {
				r.Equal(exp, intv.Containts(val))
			}
		})
	}

	t.Run("panic", func(t *testing.T) {
		require.Panics(t, func() {
			interval.New(2, 1)
		})
	})
}
