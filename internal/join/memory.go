package join

import (
	"fmt"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
	"mltwist/pkg/model"

	"github.com/zyedidia/generic/interval"
)

type memMap map[expr.Key]*memory

func (m memMap) load(key expr.Key, addr model.Addr, w expr.Width) (expr.Expr, bool) {
	mem, ok := m[key]
	if !ok {
		return nil, false
	}

	return mem.load(addr, w)
}

func (m memMap) store(key expr.Key, addr model.Addr, ex expr.Expr, w expr.Width) {
	mem, ok := m[key]
	if !ok {
		mem = newMemory()
		m[key] = mem
	}

	mem.store(addr, ex, w)
}

type memory struct {
	t *interval.Tree[model.Addr, cutExpr]
}

func newMemory() *memory {
	return &memory{t: interval.New[model.Addr, cutExpr]()}
}

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

func (m *memory) load(addr model.Addr, w expr.Width) (expr.Expr, bool) {
	end := addr + model.Addr(w)

	intervals := m.t.Overlaps(addr, end)
	if !wholeInterval(addr, end, intervals) {
		return nil, false
	}

	var finalEx expr.Expr
	if low := intervals[0].Low; low == addr {
		finalEx = intervals[0].Val.expr()
	} else {
		finalEx = intervals[0].Val.cutBegin(expr.Width(addr - low)).expr()
	}

	for _, o := range intervals[1:] {
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

func (m *memory) store(addr model.Addr, ex expr.Expr, w expr.Width) {
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

type cutExpr struct {
	ex         expr.Expr
	begin, end expr.Width
}

func (c cutExpr) width() expr.Width { return c.end - c.begin }

// cutBegin cuts bytes from an expression beginning to achieve expression of
// length length.
func (c cutExpr) cutBegin(length expr.Width) cutExpr {
	if l := c.width(); l < length {
		panic(fmt.Sprintf("bug: expr is not long enough: %d < %d", l, length))
	}

	return cutExpr{
		ex:    c.ex,
		begin: c.end - length,
		end:   c.end,
	}
}

// cutEnd cuts bytes from an expression end to achieve expression of length
// length.
func (c cutExpr) cutEnd(length expr.Width) cutExpr {
	if l := c.width(); l < length {
		panic(fmt.Sprintf("bug: expr is not long enough: %d < %d", l, length))
	}

	return cutExpr{
		ex:    c.ex,
		begin: c.begin,
		end:   c.begin + length,
	}
}

// expr returns an expression of width c.width() containing bytes [begin, end)
// from the original expression.
func (c cutExpr) expr() expr.Expr {
	if c.begin >= c.end {
		panic(fmt.Sprintf("invalid begin and end: %d >= %d", c.begin, c.end))
	}

	ex := c.ex
	if c.begin > 0 {
		shift := expr.ConstFromUint(uint16(c.begin) * 8)
		ex = expr.NewBinary(expr.Rsh, ex, shift, ex.Width())
	}

	return exprtransform.SetWidth(ex, c.end-c.begin)
}
