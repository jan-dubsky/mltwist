package interval_test

import (
	"mltwist/internal/state/interval"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap_New(t *testing.T) {
	tests := []struct {
		name string
		is   []interval.Interval[int]
		exp  []interval.Interval[int]
	}{
		{
			name: "growing_sequence",
			is: []interval.Interval[int]{
				interval.New(1, 3),
				interval.New(4, 7),
				interval.New(12, 16),
			},
			exp: []interval.Interval[int]{
				interval.New(1, 3),
				interval.New(4, 7),
				interval.New(12, 16),
			},
		}, {
			name: "growing_sequence_joining_but_front",
			is: []interval.Interval[int]{
				interval.New(0, 0),
				interval.New(1, 3),
				interval.New(3, 4),
				interval.New(4, 7),
				interval.New(7, 16),
			},
			exp: []interval.Interval[int]{
				interval.New(0, 0),
				interval.New(1, 16),
			},
		}, {
			name: "growing_sequence_joining_all",
			is: []interval.Interval[int]{
				interval.New(0, 1),
				interval.New(1, 3),
				interval.New(3, 4),
				interval.New(4, 7),
				interval.New(7, 16),
			},
			exp: []interval.Interval[int]{interval.New(0, 16)},
		}, {
			name: "growing_sequence_joining_but_end",
			is: []interval.Interval[int]{
				interval.New(0, 1),
				interval.New(1, 3),
				interval.New(3, 4),
				interval.New(4, 7),
				interval.New(7, 16),
				interval.New(17, 23),
			},
			exp: []interval.Interval[int]{
				interval.New(0, 16),
				interval.New(17, 23),
			},
		}, {
			name: "interval_overlapping_one_another",
			is: []interval.Interval[int]{
				interval.New(0, 1),
				interval.New(1, 4),
				interval.New(3, 4),
			},
			exp: []interval.Interval[int]{interval.New(0, 4)},
		}, {
			name: "intervals_unsorted",
			is: []interval.Interval[int]{
				interval.New(5, 7),
				interval.New(1, 4),
				interval.New(0, 1),
				interval.New(7, 9),
			},
			exp: []interval.Interval[int]{
				interval.New(0, 4),
				interval.New(5, 9),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			is := interval.NewMap(tt.is...)
			require.Equal(t, tt.exp, is.Intervals())
		})
	}
}

func TestMap_Add(t *testing.T) {
	tests := []struct {
		name string
		is1  interval.Map[int]
		is2  interval.Map[int]
		exp  interval.Map[int]
	}{
		{
			name: "non_overlapping",
			is1: interval.NewMap(
				interval.New(1, 2),
				interval.New(3, 5),
				interval.New(6, 7),
			),
			is2: interval.NewMap(
				interval.New(-1, 0),
				interval.New(8, 13),
			),
			exp: interval.NewMap(
				interval.New(-1, 0),
				interval.New(1, 2),
				interval.New(3, 5),
				interval.New(6, 7),
				interval.New(8, 13),
			),
		}, {
			name: "complementary",
			is1: interval.NewMap(
				interval.New(1, 2),
				interval.New(3, 5),
				interval.New(6, 7),
			),
			is2: interval.NewMap(
				interval.New(2, 3),
				interval.New(5, 6),
			),
			exp: interval.NewMap(interval.New(1, 7)),
		}, {
			name: "one_subset_of_another",
			is1: interval.NewMap(
				interval.New(1, 4),
				interval.New(5, 8),
				interval.New(9, 14),
			),
			is2: interval.NewMap(
				interval.New(2, 3),
				interval.New(5, 7),
				interval.New(11, 14),
			),
			exp: interval.NewMap(
				interval.New(1, 4),
				interval.New(5, 8),
				interval.New(9, 14),
			),
		}, {
			name: "complementary_overlapping",
			is1: interval.NewMap(
				interval.New(1, 2),
				interval.New(3, 5),
				interval.New(6, 7),
				interval.New(12, 15),
				interval.New(18, 25),
				interval.New(27, 31),
			),
			is2: interval.NewMap(
				interval.New(0, 3),
				interval.New(2, 6),
				interval.New(7, 14),
				interval.New(7, 14),
				interval.New(15, 21),
				interval.New(20, 29),
			),
			exp: interval.NewMap(interval.New(0, 31)),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			is := interval.Add(tt.is1, tt.is2)
			require.Equal(t, tt.exp, is)

			// Add(a, b) is equivalent to Add(b, a).
			is = interval.Add(tt.is2, tt.is1)
			require.Equal(t, tt.exp, is)
		})
	}
}

func TestMap_Sub(t *testing.T) {
	tests := []struct {
		name string
		is1  interval.Map[int]
		is2  interval.Map[int]
		exp  interval.Map[int]
	}{
		{
			name: "non_overlapping",
			is1: interval.NewMap(
				interval.New(1, 2),
				interval.New(3, 5),
				interval.New(6, 7),
			),
			is2: interval.NewMap(
				interval.New(-1, 0),
				interval.New(8, 13),
			),
			exp: interval.NewMap(
				interval.New(1, 2),
				interval.New(3, 5),
				interval.New(6, 7),
			),
		}, {
			name: "overlapping",
			is1: interval.NewMap(
				interval.New(1, 5),
				interval.New(6, 8),
				interval.New(12, 17),
			),
			is2: interval.NewMap(
				interval.New(4, 6),
				interval.New(6, 16),
			),
			exp: interval.NewMap(
				interval.New(1, 4),
				interval.New(16, 17),
			),
		}, {
			name: "split_interval",
			is1: interval.NewMap(
				interval.New(1, 17),
			),
			is2: interval.NewMap(
				interval.New(3, 4),
				interval.New(7, 9),
				interval.New(12, 16),
			),
			exp: interval.NewMap(
				interval.New(1, 3),
				interval.New(4, 7),
				interval.New(9, 12),
				interval.New(16, 17),
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			is := interval.Sub(tt.is1, tt.is2)
			require.Equal(t, tt.exp, is)
		})
	}
}
