package memory

import (
	"mltwist/internal/state/interval"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"

	intervaltree "github.com/zyedidia/generic/interval"
)

var _ Memory = &Sparse{}

// Sparse is a sparse memory representation.
//
// There are multiple ways how we can represent a memory address space. The
// simplest one is to use just an array of bytes. This representation would be
// as close to reality as possible, but it has several drawbacks. First of all,
// we'd need to hold the whole address space of the program in decompiler
// memory. This sounds possible for 32bit platform, but 64bit memory address
// space is simply too big to handle. Second even if we were able to store the
// whole address space, we would struggle to differentiate which bytes we wrote
// and which of them have just default (zero) value. As for some use-cases, as
// for example interpreting data flow via memory inside a single basic block, we
// need to know whether memory was ever written or not, we'd need some
// additional data on top of a byte array to support this use-case.
//
// Instead of a full byte array, we can represent memory as sparse byte array.
// This would address both problems listed above. The problem with memory
// consumption is addresses by simply storing only those blocks which were
// already written and significantly lowering memory consumption by not storing
// other bytes. Any representation of sparse memory also contains an information
// which memory was written and which was not as those blocks which were not
// written are just not included in the memory. This way, we are able to
// recognize memory addresses which were never written.
type Sparse struct {
	t *intervaltree.Tree[model.Addr, cutExpr]
}

func NewSparse() *Sparse {
	return &Sparse{t: intervaltree.New[model.Addr, cutExpr]()}
}

// wholeInterval checks that a list of intervals ints spans the whole [begin,
// end) interval.
func wholeInterval(begin, end model.Addr, ints []intervaltree.KV[model.Addr, cutExpr]) bool {
	if len(ints) == 0 || ints[0].Low > begin {
		return false
	}

	lastEnd := ints[0].High
	for _, o := range ints[1:] {
		if o.Low != lastEnd {
			return false
		}
		lastEnd = o.High
	}

	if lastEnd < end {
		return false
	}

	return true
}

func (m *Sparse) Load(addr model.Addr, w expr.Width) (expr.Expr, bool) {
	end := addr + model.Addr(w)

	ints := m.t.Overlaps(addr, end)
	if !wholeInterval(addr, end, ints) {
		return nil, false
	}

	var finalEx expr.Expr
	if low := ints[0].Low; low == addr {
		finalEx = ints[0].Val.expr()
	} else {
		finalEx = ints[0].Val.cutBegin(expr.Width(addr - low)).expr()
	}

	for _, o := range ints[1:] {
		var ex expr.Expr
		if o.High <= end {
			ex = o.Val.expr()
		} else {
			ex = o.Val.cutEnd(expr.Width(o.High - end)).expr()
		}

		ex = expr.NewBinary(expr.Lsh, ex, expr.ConstFromUint((o.Low-addr)*8), w)
		finalEx = expr.NewBinary(expr.Or, finalEx, ex, w)
	}

	return finalEx, true
}

func (m *Sparse) Store(addr model.Addr, ex expr.Expr, w expr.Width) {
	end := addr + model.Addr(w)

	overlaps := m.t.Overlaps(addr, end)
	// First remove all overlapping interval so that we are sure that all
	// Adds below will succeed.
	for _, o := range overlaps {
		m.t.Remove(o.Low)
	}

	for _, o := range overlaps {
		// Those are fully rewritten by this write.
		if addr <= o.Low && o.High <= end {
			continue
		}

		if o.Low < addr {
			m.t.Add(o.Low, addr, o.Val.cutEnd(expr.Width(addr-o.Low)))
		}
		// Note: Else branch cannot be used here as we might be in a
		// following situation: o.Low < begin < end < o.High.
		if end < o.High {
			m.t.Put(end, o.High, o.Val.cutBegin(expr.Width(o.High-end)))
		}
	}

	c := cutExpr{
		ex:    ex,
		begin: 0,
		end:   w,
	}
	m.t.Add(addr, addr+model.Addr(w), c)
}

func (m *Sparse) Missing(addr model.Addr, w expr.Width) interval.Map[model.Addr] {
	end := addr + model.Addr(w)

	ints := m.t.Overlaps(addr, end)
	if len(ints) == 0 {
		return interval.NewMap(interval.New(addr, end))
	}

	var intervals []interval.Interval[model.Addr]
	if low := ints[0].Low; addr < low {
		intervals = append(intervals, interval.New(addr, low))
	}

	lastEnd := ints[0].High
	for _, o := range ints[1:] {
		if o.Low != lastEnd {
			intervals = append(intervals, interval.New(lastEnd, o.Low))
		}
		lastEnd = o.High
	}

	if lastEnd < end {
		intervals = append(intervals, interval.New(lastEnd, end))
	}

	return interval.NewMap(intervals...)
}

func (m *Sparse) Blocks() interval.Map[model.Addr] {
	blocks := make([]interval.Interval[model.Addr], 0, m.t.Size())
	m.t.Each(func(begin, end model.Addr, val cutExpr) {
		blocks = append(blocks, interval.New(begin, end))
	})

	return interval.NewMap(blocks...)
}
