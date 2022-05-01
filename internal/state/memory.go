package state

import (
	"fmt"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"

	"github.com/zyedidia/generic/interval"
)

// MamMap is set of individual memory address spaces identified by a key.
type MemMap map[expr.Key]*Memory

// Load loads w bytes from an address addr in address space identified by key.
func (m MemMap) Load(key expr.Key, addr model.Addr, w expr.Width) (expr.Expr, bool) {
	mem, ok := m[key]
	if !ok {
		return nil, false
	}

	return mem.Load(addr, w)
}

// Store stores w bytes of ex to address addr in address space identified by
// key. If address space for key key doesn't exist yet, it's created.
func (m MemMap) Store(key expr.Key, addr model.Addr, ex expr.Expr, w expr.Width) {
	mem, ok := m[key]
	if !ok {
		mem = NewMemory()
		m[key] = mem
	}

	mem.Store(addr, ex, w)
}

func (m MemMap) Missing(key expr.Key, addr model.Addr, w expr.Width) []Interval {
	mem, ok := m[key]
	if !ok {
		return []Interval{newInterval(addr, addr+model.Addr(w))}
	}

	return mem.Missing(addr, w)
}

// Memory represents a single linear address space where expressions of
// different width are stored to and loaded from.
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
// written are just not included in the memory. Consequently we are able to
// recognize memory addresses which were never written.
//
// The last challenge of our memory design in the fact that we don't store
// bytes, but we store expressions, which might not be constant. Consequently
// there exists no way how to store those expressions as bytes as we simply
// cannot deduce their byte value. What we can do is to store raw expressions.
// Storing expression instead of bytes brings some challenges in mapping the
// expression model to a byte-oriented model. Namely a new write can rewrite
// just part of an expression written before. Consequently, we have to store an
// information which bytes of an expression are still valid and which of them
// were rewritten. On read, we are able to compose any expression as combination
// of written expression modified by bit shifts, ANDs and OSs.
//
// The logic described above might result in duplication of expression. This
// happens if a new write rewrites just some bytes in the middle of a previously
// written expression. This results in 2 places in memory being calculated by
// the same expression which differ only by bytes taken. Given that every write
// can result in duplication of at most one expression, the ultimate growth of
// expression complexity is only linear in number of writes. Consequently the
// overall number of expression used to represent the state of memory after n
// writes is always O(n) independently on the fact whether expression splitting
// happens or not. Consequently an expression splitting is not a problem as it
// doesn't increase number of expressions to evaluate significantly.
type Memory struct {
	t *interval.Tree[model.Addr, cutExpr]
}

func NewMemory() *Memory {
	return &Memory{t: interval.New[model.Addr, cutExpr]()}
}

// wholeInterval checks that a list of intervals ints spans the whole [begin,
// end) interval.
func wholeInterval(begin, end model.Addr, ints []interval.KV[model.Addr, cutExpr]) bool {
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

// Load loads w bytes from an address addr from the memory.
func (m *Memory) Load(addr model.Addr, w expr.Width) (expr.Expr, bool) {
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

		ex = expr.NewBinary(expr.Lsh, ex, expr.ConstFromUint(o.Low), w)
		finalEx = expr.NewBinary(expr.Or, finalEx, ex, w)
	}

	return finalEx, true
}

// Store stores expression ex of width w to address addr in the memory.
func (m *Memory) Store(addr model.Addr, ex expr.Expr, w expr.Width) {
	end := addr + model.Addr(w)

	overlaps := m.t.Overlaps(addr, end)
	// First remove all overlapping interval so that we are sure that all
	// Adds will succeed.
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
		// Note: else expression cannot be used here as we might be in a
		// following situation: o.Low < begin < end < o.High.
		if end < o.High {
			m.t.Put(end, o.High, o.Val.cutBegin(expr.Width(o.High-addr)))
		}
	}

	c := cutExpr{
		ex:    ex,
		begin: 0,
		end:   w,
	}
	m.t.Add(addr, addr+model.Addr(w), c)
}

type Interval struct {
	begin model.Addr
	end   model.Addr
}

func newInterval(begin, end model.Addr) Interval {
	if w := end - begin; w != model.Addr(expr.Width(w)) {
		panic(fmt.Sprintf("interval is too wide: %d", w))
	}

	return Interval{
		begin: begin,
		end:   end,
	}
}

func (i Interval) Begin() model.Addr { return i.begin }
func (i Interval) End() model.Addr   { return i.end }
func (i Interval) Width() expr.Width { return expr.Width(i.end - i.begin) }

// Missing return list of address intervals which are missing in memory to be
// able to perform full load of width w from address addr.
func (m *Memory) Missing(addr model.Addr, w expr.Width) []Interval {
	end := addr + model.Addr(w)

	ints := m.t.Overlaps(addr, end)
	if len(ints) == 0 {
		return []Interval{newInterval(addr, end)}
	}

	var intervals []Interval
	if low := ints[0].Low; addr < low {
		intervals = append(intervals, newInterval(addr, low))
	}

	lastEnd := ints[0].High
	for _, o := range ints[1:] {
		if o.Low != lastEnd {
			intervals = append(intervals, newInterval(lastEnd, o.Low))
		}
		lastEnd = o.High
	}

	if lastEnd < end {
		intervals = append(intervals, newInterval(lastEnd, end))
	}

	return intervals
}
