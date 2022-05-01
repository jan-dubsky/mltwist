package state

import (
	"fmt"
	"mltwist/internal/exprtransform"
	"mltwist/pkg/expr"
)

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
