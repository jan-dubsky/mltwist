package memory

import (
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type testByteBlock struct {
	b  model.Addr
	bs []byte
}

func (b testByteBlock) Begin() model.Addr { return b.b }
func (b testByteBlock) Bytes() []byte     { return b.bs }

func TestBytes_New(t *testing.T) {
	tests := []struct {
		name   string
		blocks []testByteBlock
		hasErr bool
		exp    []testByteBlock
	}{{
		name: "no_blocks",
	}, {
		name: "single_block",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		},
		exp: []testByteBlock{{
			b:  25,
			bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}},
	}, {
		name: "two_adjacent_blocks",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			{b: 34, bs: []byte{11, 12, 13, 14}},
		},
		exp: []testByteBlock{{
			b:  25,
			bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14},
		}},
	}, {
		name: "two_block_sequences",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			{b: 34, bs: []byte{11, 12, 13, 14}},
			{b: 60, bs: []byte{21, 22, 23, 24, 25, 26, 27}},
		},
		exp: []testByteBlock{{
			b:  25,
			bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14},
		}, {
			b:  60,
			bs: []byte{21, 22, 23, 24, 25, 26, 27},
		}},
	}, {
		name: "overlapping_blocks",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			{b: 24, bs: []byte{11, 12, 13, 14}},
		},
		hasErr: true,
	}, {
		name: "multiblock_recombination",
		blocks: []testByteBlock{
			{b: 38, bs: []byte{0xfe, 0xff}},
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			{b: 34, bs: []byte{11, 12, 13, 14}},
			{b: 40, bs: []byte{21, 22, 23, 24, 25, 26, 27}},
		},
		exp: []testByteBlock{{
			b: 25,
			bs: []byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14,
				0xfe, 0xff, 21, 22, 23, 24, 25, 26, 27,
			},
		}},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			blocks := make([]ByteBlock, len(tt.blocks))
			for i, b := range tt.blocks {
				blocks[i] = b
			}

			mem, err := NewBytes(blocks)
			if tt.hasErr {
				r.Error(err)
				return
			}

			r.Equal(len(tt.exp), len(mem.blocks))
			for i, e := range tt.exp {
				r.Equal(e.b, mem.blocks[i].begin)
				r.Equal(e.bs, mem.blocks[i].bytes)
			}

		})
	}
}

func TestBytes_Load(t *testing.T) {
	tests := []struct {
		name   string
		blocks []testByteBlock
		addr   model.Addr
		w      expr.Width
		exp    []byte
	}{{
		name: "empty_memory",
		addr: 25,
		w:    4,
	}, {
		name: "block_start",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		},
		addr: 25,
		w:    8,
		exp:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
	}, {
		name: "block_middle",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 44,
		w:    2,
		exp:  []byte{15, 16},
	}, {
		name: "block_end",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 46,
		w:    4,
		exp:  []byte{17, 18, 19, 20},
	}, {
		name: "block_behind_end",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 47,
		w:    4,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			blocks := make([]ByteBlock, len(tt.blocks))
			for i, b := range tt.blocks {
				blocks[i] = b
			}

			mem, err := NewBytes(blocks)
			r.NoError(err)

			ex, ok := mem.Load(tt.addr, tt.w)
			if tt.exp == nil {
				r.False(ok)
				r.Nil(ex)
				return
			}

			r.True(ok)
			r.IsType(expr.Const{}, ex)
			r.Equal(tt.exp, ex.(expr.Const).Bytes())
		})
	}
}

func TestBytes_Store(t *testing.T) {
	tests := []struct {
		name   string
		blocks []testByteBlock
		addr   model.Addr
		ex     expr.Const
		w      expr.Width
		exp    []testByteBlock
	}{{
		name: "write_into_block",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 27,
		ex:   expr.NewConst([]byte{0x80, 0x81, 0x82, 0x83}, expr.Width32),
		w:    expr.Width32,
		exp: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 0x80, 0x81, 0x82, 0x83, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
	}, {
		name: "narrower_write",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 27,
		ex:   expr.NewConst([]byte{0x80, 0x81, 0x82, 0x83}, expr.Width32),
		w:    expr.Width16,
		exp: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 0x80, 0x81, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
	}, {
		name: "append_to_block",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 35,
		ex:   expr.NewConst([]byte{0x80, 0x81, 0x82, 0x83}, expr.Width32),
		w:    expr.Width32,
		exp: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
				0x80, 0x81, 0x82, 0x83}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
	}, {
		name: "append_to_block_overlapping",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 33,
		ex:   expr.NewConst([]byte{0x80, 0x81, 0x82, 0x83}, expr.Width32),
		w:    expr.Width32,
		exp: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 0x80, 0x81, 0x82, 0x83}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
	}, {
		name: "prepend_to_block",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 36,
		ex:   expr.NewConst([]byte{0x80, 0x81, 0x82, 0x83}, expr.Width32),
		w:    expr.Width32,
		exp: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 36, bs: []byte{0x80, 0x81, 0x82, 0x83, 11, 12, 13, 14,
				15, 16, 17, 18, 19, 20}},
		},
	}, {
		name: "prepend_to_block_overlapping",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 38,
		ex:   expr.NewConst([]byte{0x80, 0x81, 0x82, 0x83}, expr.Width32),
		w:    expr.Width32,
		exp: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 38, bs: []byte{0x80, 0x81, 0x82, 0x83, 13, 14,
				15, 16, 17, 18, 19, 20}},
		},
	}, {
		name: "join_blocks",
		blocks: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{b: 40, bs: []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
		addr: 34,
		ex: expr.NewConst(
			[]byte{0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87},
			expr.Width64,
		),
		w: expr.Width64,
		exp: []testByteBlock{
			{b: 25, bs: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0x80,
				0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 13,
				14, 15, 16, 17, 18, 19, 20,
			}},
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			blocks := make([]ByteBlock, len(tt.blocks))
			for i, b := range tt.blocks {
				blocks[i] = b
			}

			mem, err := NewBytes(blocks)
			r.NoError(err)

			mem.Store(tt.addr, tt.ex, tt.w)

			r.Equal(len(tt.exp), len(mem.blocks))
			for i, e := range tt.exp {
				r.Equal(e.b, mem.blocks[i].begin)
				r.Equal(e.bs, mem.blocks[i].bytes)
			}
		})
	}

	t.Run("panic", func(t *testing.T) {
		r := require.New(t)

		mem, err := NewBytes(nil)
		r.NoError(err)
		r.Panics(func() {
			mem.Store(55, expr.NewRegLoad("r1", expr.Width8), expr.Width32)
		})
	})
}
