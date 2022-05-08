package state

import (
	"mltwist/internal/state/interval"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"

	intervaltree "github.com/zyedidia/generic/interval"
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

func (m MemMap) Missing(key expr.Key, addr model.Addr, w expr.Width) interval.Map[model.Addr] {
	mem, ok := m[key]
	if !ok {
		return interval.NewMap(interval.New(addr, addr+model.Addr(w)))
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
// written are just not included in the memory. This way, we are able to
// recognize memory addresses which were never written.
//
// The last challenge of our memory design in the fact that we don't store
// bytes, but we store expressions, which might not be constant. Consequently
// there exists no way how to store those expressions as bytes as we simply
// cannot deduce their byte value. What we can do is to store raw expressions.
// Storing expression instead of bytes brings some challenges in mapping the
// expression model to a byte-oriented model. Namely a new write can rewrite
// just part of an expression written before. For this reason, we have to store
// an information which bytes of an expression are still valid and which of them
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
// happens or not. Given this, the expression splitting is not much of an issue
// as it doesn't increase number of expressions to evaluate significantly.
type Memory struct {
	t *intervaltree.Tree[model.Addr, cutExpr]
}

func NewMemory() *Memory {
	return &Memory{t: intervaltree.New[model.Addr, cutExpr]()}
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

		ex = expr.NewBinary(expr.Lsh, ex, expr.ConstFromUint((o.Low-addr)*8), w)
		finalEx = expr.NewBinary(expr.Or, finalEx, ex, w)
	}

	return finalEx, true
}

// Store stores expression ex of width w to address addr in the memory.
func (m *Memory) Store(addr model.Addr, ex expr.Expr, w expr.Width) {
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

// Missing return list of address intervals which are missing in memory to be
// able to perform full load of width w from address addr.
func (m *Memory) Missing(addr model.Addr, w expr.Width) interval.Map[model.Addr] {
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

// Blocks returns a list of continuous blocks stored in the memory.
func (m *Memory) Blocks() interval.Map[model.Addr] {
	blocks := make([]interval.Interval[model.Addr], 0, m.t.Size())
	m.t.Each(func(begin, end model.Addr, val cutExpr) {
		blocks = append(blocks, interval.New(begin, end))
	})

	return interval.NewMap(blocks...)
}
