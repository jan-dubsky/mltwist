package cursor

import (
	"decomp/internal/console/internal/lines"
	"fmt"
)

type Cursor struct {
	l *lines.Lines

	value int
}

func New(l *lines.Lines) *Cursor {
	return &Cursor{
		l:     l,
		value: 0,
	}
}

func (c Cursor) Value() int { return c.value }
func (c *Cursor) Set(o int) error {
	if err := c.checkOffset(o); err != nil {
		return fmt.Errorf("new offset value is invalid: %w", err)
	}

	c.value = o
	return nil
}

func (c Cursor) checkOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset cannot be negative: %d", offset)
	}
	if l := c.l.Len(); offset >= l {
		return fmt.Errorf("offset is too high: %d > %d", offset, l)
	}

	return nil
}
