package deps

import (
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockMovable struct {
	id     uint
	length model.Addr

	idx  int
	addr model.Addr
}

func newMock(id uint, length model.Addr) *mockMovable {
	return &mockMovable{
		id:     id,
		length: length,
	}
}

func (m *mockMovable) setIndex(i int)       { m.idx = i }
func (m *mockMovable) setAddr(a model.Addr) { m.addr = a }
func (m mockMovable) Addr() model.Addr      { return m.addr }
func (m mockMovable) NextAddr() model.Addr  { return m.addr + m.length }

func TestMove(t *testing.T) {
	tests := []struct {
		name     string
		ms       []*mockMovable
		from     int
		to       int
		expected []uint
	}{
		{
			name: "fwd_equidistant",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 2),
				newMock(2, 2),
				newMock(3, 2),
				newMock(4, 2),
			},
			from:     1,
			to:       4,
			expected: []uint{0, 2, 3, 4, 1},
		}, {
			name: "back_equidistant",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 2),
				newMock(2, 2),
				newMock(3, 2),
				newMock(4, 2),
			},
			from:     3,
			to:       0,
			expected: []uint{3, 0, 1, 2, 4},
		}, {
			name: "fwd_equidistant_last",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 2),
				newMock(2, 2),
				newMock(3, 2),
				newMock(4, 2),
			},
			from:     0,
			to:       2,
			expected: []uint{1, 2, 0, 3, 4},
		}, {
			name: "back_equidistant_last",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 2),
				newMock(2, 2),
				newMock(3, 2),
				newMock(4, 2),
			},
			from:     4,
			to:       2,
			expected: []uint{0, 1, 4, 2, 3},
		}, {
			name: "fwd_different",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 3),
				newMock(2, 5),
				newMock(3, 7),
				newMock(4, 11),
				newMock(5, 13),
			},
			from:     2,
			to:       5,
			expected: []uint{0, 1, 3, 4, 5, 2},
		}, {
			name: "back_different",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 3),
				newMock(2, 5),
				newMock(3, 7),
				newMock(4, 11),
				newMock(5, 13),
			},
			from:     4,
			to:       1,
			expected: []uint{0, 4, 1, 2, 3, 5},
		}, {
			name: "fwd_different_last",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 3),
				newMock(2, 5),
				newMock(3, 7),
				newMock(4, 11),
				newMock(5, 13),
			},
			from:     0,
			to:       4,
			expected: []uint{1, 2, 3, 4, 0, 5},
		}, {
			name: "back_different_last",
			ms: []*mockMovable{
				newMock(0, 2),
				newMock(1, 3),
				newMock(2, 5),
				newMock(3, 7),
				newMock(4, 11),
				newMock(5, 13),
			},
			from:     5,
			to:       3,
			expected: []uint{0, 1, 2, 5, 3, 4},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var addr model.Addr
			for i := range tt.ms {
				tt.ms[i].idx = i
				tt.ms[i].addr = addr
				addr += tt.ms[i].length
			}

			move(tt.ms, tt.from, tt.to)

			r := require.New(t)
			for i, e := range tt.expected {
				t.Logf("expected index: %d", i)
				r.Equal(e, tt.ms[i].id)
			}

			addr = 0
			for i, m := range tt.ms {
				t.Logf("movable index: %d", i)
				r.Equal(i, m.idx)
				r.Equal(addr, m.addr)
				addr += m.length
			}
		})
	}
}
