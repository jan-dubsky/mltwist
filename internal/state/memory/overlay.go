package memory

import (
	"mltwist/internal/state/interval"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"
	"sort"
)

var _ Memory = &Overlay{}

// Overlay is an implementation of Memory interface which composes 2 Memory
// layers into a single memory.
//
// All writes to an Overlay memory are written in the overlay memory layer. On
// the other hand all other methods of this memory type behave as if overlay
// layer was put on top of the base layer. For example memory loads first load
// from overlay layer but if not all bytes to be read are present in the overlay
// layer, they are read from the base layer.
//
// As loads can span multiple bytes, it's possible that some bytes are present
// in the overlay layer and other are present in the base layer. In such a case,
// the load result is combination of expressions from those two layers using bit
// shifts, ANDs and ORs.
//
// The design of this memory type is strongly inspired by overlayfs in Linux
// (used for example by docker) which allows to compose multiple layers of
// memory on top of one another. This memory type behaves in the same way and it
// allows us to compose multiple different memories into the final memory space.
//
// Overlay layer can be technically used recursively to construct multi-layer
// memory. /even though such usage is possible it's strongly discouraged. The
// reasoning behind is that every Overlay layer adds a non-trivial computational
// and allocation complexity to every memory load. In case of many overlay
// layers being stack one on another, this effect will sum up. This is by the
// way one of the major criticisms for overlayfs in Linux as well.
//
// As the base layer of Overlay layer is always readonly, it's strongly
// recommended to replace recursive usage of Overlay layer by a single memory
// layer storing Load results for all memory layer in the stack. In other words,
// it's recommended to squash multiple layers on top of one another into a
// single layer to avoid the overhead of many overlay layers on top of one
// another.
type Overlay struct {
	base    Memory
	overlay Memory
}

func NewOverlay(base Memory, overlay Memory) *Overlay {
	return &Overlay{
		base:    base,
		overlay: overlay,
	}
}

// Base returns the base layer of o.
func (o Overlay) Base() Memory { return o.base }

// Overlay returns the overlay layer of o.
func (o Overlay) Overlay() Memory { return o.overlay }

type rangeRead struct {
	intv interval.Interval[model.Addr]
	ex   expr.Expr
}

func offsetExpr(ex expr.Expr, bytes expr.Width, w expr.Width) expr.Expr {
	return expr.NewBinary(expr.Lsh, ex, expr.ConstFromUint(bytes.Bits()), w)
}

func (o *Overlay) Load(addr model.Addr, w expr.Width) (expr.Expr, bool) {
	missing := o.overlay.Missing(addr, w)
	if missing.Len() == 0 {
		return o.overlay.Load(addr, w)
	}

	wholeRange := interval.NewMap(interval.New(addr, addr+model.Addr(w)))
	if wholeRange.Equal(missing) {
		return o.base.Load(addr, w)
	}

	overlay := interval.MapIntersect(wholeRange, missing)

	reads := make([]rangeRead, 0, missing.Len()+overlay.Len())
	for _, intv := range missing.Intervals() {
		ex, ok := o.base.Load(intv.Begin(), expr.Width(intv.End()))
		if !ok {
			return nil, false
		}

		reads = append(reads, rangeRead{intv: intv, ex: ex})
	}
	for _, intv := range overlay.Intervals() {
		ex, ok := o.overlay.Load(intv.Begin(), expr.Width(intv.End()))
		if !ok {
			return nil, false
		}

		reads = append(reads, rangeRead{intv: intv, ex: ex})
	}

	sort.Slice(reads, func(i, j int) bool {
		return reads[i].intv.Begin() < reads[j].intv.Begin()
	})

	finalEx := reads[0].ex
	for _, r := range reads[1:] {
		ex := offsetExpr(r.ex, expr.Width(r.intv.Begin()), w)
		finalEx = expr.NewBinary(expr.Or, finalEx, ex, w)
	}

	return finalEx, true
}

func (o *Overlay) Store(addr model.Addr, ex expr.Expr, w expr.Width) {
	o.overlay.Store(addr, ex, w)
}

func (o *Overlay) Missing(addr model.Addr, w expr.Width) interval.Map[model.Addr] {
	return interval.MapIntersect(o.base.Missing(addr, w), o.overlay.Missing(addr, w))
}

func (o *Overlay) Blocks() interval.Map[model.Addr] {
	return interval.MapUnion(o.base.Blocks(), o.overlay.Blocks())
}
