package memory_test

import (
	"bytes"
	"mltwist/internal/exprtransform"
	"mltwist/internal/state/interval"
	"mltwist/internal/state/memory"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type mem struct {
	sig   byte
	intvs interval.Map[model.Addr]
}

func (m *mem) Load(addr model.Addr, w expr.Width) (expr.Expr, bool) {
	if m.Missing(addr, w).Len() != 0 {
		return nil, false
	}

	bs := bytes.Repeat([]byte{m.sig}, int(w))
	return expr.NewConst(bs, w), true
}

func (m *mem) Store(_ model.Addr, _ expr.Expr, _expr.Width) {}

func (m *mem) Missing(addr model.Addr, w expr.Width) interval.Map[model.Addr] {
	intv := interval.NewMap(interval.New(addr, addr+model.Addr(w)))
	return interval.MapComplement(intv, m.intvs)
}

func (m *mem) Blocks() interval.Map[model.Addr] { return m.intvs }

func TestOverlay_Load(t *testing.T) {
	base := &mem{
		sig: 0x11,
		intvs: interval.NewMap(
			interval.New[model.Addr](1, 5),
			interval.New[model.Addr](6, 24),
			interval.New[model.Addr](32, 64),
			interval.New[model.Addr](68, 72),
		),
	}

	overlay := &mem{
		sig: 0xff,
		intvs: interval.NewMap(
			interval.New[model.Addr](5, 6),
			interval.New[model.Addr](24, 30),
			interval.New[model.Addr](62, 70),
		),
	}

	tests := []struct {
		name string
		addr model.Addr
		w    expr.Width
		exp  []byte
	}{{
		name: "load_from_base",
		addr: 1,
		w:    expr.Width32,
		exp:  []byte{0x11, 0x11, 0x11, 0x11},
	}, {
		name: "overlay_bridging_base",
		addr: 4,
		w:    expr.Width32,
		exp:  []byte{0x11, 0xff, 0x11, 0x11},
	}, {
		name: "one_following_another",
		addr: 22,
		w:    expr.Width64,
		exp:  []byte{0x11, 0x11, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}, {
		name: "one_missing_byte",
		addr: 23,
		w:    expr.Width64,
	}, {
		name: "no_data",
		addr: 30,
		w:    expr.Width16,
		exp:  nil,
	}, {
		name: "overlay_over_base",
		addr: 60,
		w:    expr.Width64,
		exp:  []byte{0x11, 0x11, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}, {
		name: "full_overlay",
		addr: 62,
		w:    expr.Width64,
		exp:  []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			mem := memory.NewOverlay(base, overlay)

			ex, ok := mem.Load(tt.addr, tt.w)
			if tt.exp == nil {
				r.False(ok)
				return
			}

			r.True(ok)

			c := exprtransform.ConstFold(ex).(expr.Const)
			r.Equal(tt.exp, c.Bytes())
		})
	}

	require.Equal(t, 0, len(base.writes))
}
