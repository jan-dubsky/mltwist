package state

import (
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type testInterval struct {
	begin    model.Addr
	end      model.Addr
	ex       expr.Expr
	cutBegin expr.Width
	cutEnd   expr.Width
}

func TestMemory_Store(t *testing.T) {
	tests := []struct {
		name string
		addr model.Addr
		e    expr.Expr
		w    expr.Width
		exp  []testInterval
	}{{
		name: "first_write",
		addr: 66,
		e:    expr.NewBinary(expr.Add, expr.One, expr.Zero, expr.Width64),
		w:    expr.Width16,
		exp: []testInterval{{
			begin:    66,
			end:      68,
			ex:       expr.NewBinary(expr.Add, expr.One, expr.Zero, expr.Width64),
			cutBegin: 0,
			cutEnd:   2,
		}},
	}, {
		name: "append",
		addr: 68,
		e:    expr.ConstFromUint[uint8](68),
		w:    expr.Width32,
		exp: []testInterval{{
			begin:    66,
			end:      68,
			ex:       expr.NewBinary(expr.Add, expr.One, expr.Zero, expr.Width64),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    68,
			end:      72,
			ex:       expr.ConstFromUint[uint8](68),
			cutBegin: 0,
			cutEnd:   4,
		}},
	}, {
		name: "prepend",
		addr: 58,
		e:    expr.ConstFromUint[uint64](0xabcd12345678abcd),
		w:    expr.Width64,
		exp: []testInterval{{
			begin:    58,
			end:      66,
			ex:       expr.ConstFromUint[uint64](0xabcd12345678abcd),
			cutBegin: 0,
			cutEnd:   8,
		}, {
			begin:    66,
			end:      68,
			ex:       expr.NewBinary(expr.Add, expr.One, expr.Zero, expr.Width64),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    68,
			end:      72,
			ex:       expr.ConstFromUint[uint8](68),
			cutBegin: 0,
			cutEnd:   4,
		}},
	}, {
		name: "overlap_end",
		addr: 64,
		e:    expr.ConstFromInt[int16](-1),
		w:    expr.Width16,
		exp: []testInterval{{
			begin:    58,
			end:      64,
			ex:       expr.ConstFromUint[uint64](0xabcd12345678abcd),
			cutBegin: 0,
			cutEnd:   6,
		}, {
			begin:    64,
			end:      66,
			ex:       expr.ConstFromInt[int16](-1),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    66,
			end:      68,
			ex:       expr.NewBinary(expr.Add, expr.One, expr.Zero, expr.Width64),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    68,
			end:      72,
			ex:       expr.ConstFromUint[uint8](68),
			cutBegin: 0,
			cutEnd:   4,
		}},
	}, {
		name: "overlap_begin",
		addr: 58,
		e:    expr.ConstFromUint[uint32](0xfedcfedc),
		w:    expr.Width32,
		exp: []testInterval{{
			begin:    58,
			end:      62,
			ex:       expr.ConstFromUint[uint32](0xfedcfedc),
			cutBegin: 0,
			cutEnd:   4,
		}, {
			begin:    62,
			end:      64,
			ex:       expr.ConstFromUint[uint64](0xabcd12345678abcd),
			cutBegin: 4,
			cutEnd:   6,
		}, {
			begin:    64,
			end:      66,
			ex:       expr.ConstFromInt[int16](-1),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    66,
			end:      68,
			ex:       expr.NewBinary(expr.Add, expr.One, expr.Zero, expr.Width64),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    68,
			end:      72,
			ex:       expr.ConstFromUint[uint8](68),
			cutBegin: 0,
			cutEnd:   4,
		}},
	}, {
		name: "overlap_front_and_end",
		addr: 61,
		e:    expr.ConstFromUint[uint16](0),
		w:    expr.Width16,
		exp: []testInterval{{
			begin:    58,
			end:      61,
			ex:       expr.ConstFromUint[uint32](0xfedcfedc),
			cutBegin: 0,
			cutEnd:   3,
		}, {
			begin:    61,
			end:      63,
			ex:       expr.ConstFromUint[uint16](0),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    63,
			end:      64,
			ex:       expr.ConstFromUint[uint64](0xabcd12345678abcd),
			cutBegin: 5,
			cutEnd:   6,
		}, {
			begin:    64,
			end:      66,
			ex:       expr.ConstFromInt[int16](-1),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    66,
			end:      68,
			ex:       expr.NewBinary(expr.Add, expr.One, expr.Zero, expr.Width64),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    68,
			end:      72,
			ex:       expr.ConstFromUint[uint8](68),
			cutBegin: 0,
			cutEnd:   4,
		}},
	}, {
		name: "overwrite_intervals_fully",
		addr: 62,
		e:    expr.ConstFromUint[uint64](0x7788778877887788),
		w:    expr.Width64,
		exp: []testInterval{{
			begin:    58,
			end:      61,
			ex:       expr.ConstFromUint[uint32](0xfedcfedc),
			cutBegin: 0,
			cutEnd:   3,
		}, {
			begin:    61,
			end:      62,
			ex:       expr.ConstFromUint[uint16](0),
			cutBegin: 0,
			cutEnd:   1,
		}, {
			begin:    62,
			end:      70,
			ex:       expr.ConstFromUint[uint64](0x7788778877887788),
			cutBegin: 0,
			cutEnd:   8,
		}, {
			begin:    70,
			end:      72,
			ex:       expr.ConstFromUint[uint8](68),
			cutBegin: 2,
			cutEnd:   4,
		}},
	}}

	mem := NewMemory()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			mem.Store(tt.addr, tt.e, tt.w)

			var last model.Addr
			var i int
			mem.t.Each(func(begin, end model.Addr, val cutExpr) {
				t.Logf("entry: %d", i)

				r.LessOrEqual(last, begin)
				last = end

				r.Equal(tt.exp[i].begin, begin)
				r.Equal(tt.exp[i].end, end)

				c := cutExpr{
					ex:    tt.exp[i].ex,
					begin: tt.exp[i].cutBegin,
					end:   tt.exp[i].cutEnd,
				}
				r.Equal(c, val)

				i++
			})
		})
	}
}

func TestMemory_Load(t *testing.T) {
	tests := []struct {
		name  string
		state []testInterval
		addr  model.Addr
		w     expr.Width
		exp   expr.Expr
	}{{
		name: "empty_memory",
		addr: 68,
		w:    expr.Width32,
	}, {
		name: "missing_read_end",
		state: []testInterval{{
			begin:    68,
			end:      70,
			ex:       expr.ConstFromUint[uint8](0),
			cutBegin: 0,
			cutEnd:   2,
		}},
		addr: 68,
		w:    expr.Width32,
	}, {
		name: "missing_read_begin",
		state: []testInterval{{
			begin:    68,
			end:      70,
			ex:       expr.ConstFromUint[uint8](0),
			cutBegin: 0,
			cutEnd:   2,
		}},
		addr: 66,
		w:    expr.Width32,
	}, {
		name: "single_read",
		state: []testInterval{{
			begin:    68,
			end:      70,
			ex:       expr.ConstFromUint[uint8](0),
			cutBegin: 0,
			cutEnd:   2,
		}},
		addr: 68,
		w:    expr.Width16,
		exp:  expr.NewConstUint[uint8](0, expr.Width16),
	}, {
		name: "read_2_exprs",
		state: []testInterval{{
			begin:    68,
			end:      70,
			ex:       expr.ConstFromUint[uint8](0),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    70,
			end:      72,
			ex:       expr.ConstFromInt[int16](-1),
			cutBegin: 0,
			cutEnd:   2,
		}},
		addr: 68,
		w:    expr.Width32,
		exp:  expr.ConstFromUint[uint32](0xffff0000),
	}, {
		name: "read_with_overlaps",
		state: []testInterval{{
			begin:    68,
			end:      70,
			ex:       expr.ConstFromUint[uint8](0),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    70,
			end:      72,
			ex:       expr.ConstFromInt[int16](-1),
			cutBegin: 0,
			cutEnd:   2,
		}, {
			begin:    72,
			end:      76,
			ex:       expr.ConstFromUint[uint32](0x44444444),
			cutBegin: 0,
			cutEnd:   4,
		}, {
			begin:    76,
			end:      78,
			ex:       expr.ConstFromUint[uint16](0x9988),
			cutBegin: 0,
			cutEnd:   2,
		}},
		addr: 69,
		w:    expr.Width64,
		exp:  expr.ConstFromUint[uint64](0x8844444444ffff00),
	}, {
		name: "missing_byte_in_the middle",
		state: []testInterval{{
			begin:    68,
			end:      72,
			ex:       expr.ConstFromUint[uint8](0),
			cutBegin: 0,
			cutEnd:   4,
		}, {
			begin:    73,
			end:      77,
			ex:       expr.ConstFromUint[uint32](0x44444444),
			cutBegin: 0,
			cutEnd:   4,
		}, {
			begin:    77,
			end:      79,
			ex:       expr.ConstFromUint[uint16](0x8888),
			cutBegin: 0,
			cutEnd:   2,
		}},
		addr: 69,
		w:    expr.Width64,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			mem := NewMemory()
			for _, i := range tt.state {
				c := cutExpr{
					ex:    i.ex,
					begin: i.cutBegin,
					end:   i.cutEnd,
				}
				mem.t.Add(i.begin, i.end, c)
			}

			ex, ok := mem.Load(tt.addr, tt.w)
			if tt.exp == nil {
				r.False(ok)
				return
			}

			r.True(ok)
			r.Equal(tt.exp, exprtransform.ConstFold(ex))

		})
	}
}
