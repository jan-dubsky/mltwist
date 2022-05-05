package cursor

import (
	"fmt"
)

type Cursor struct {
	maxValue int
	value    int
}

func New(maxOffset int) *Cursor {
	return &Cursor{
		maxValue: maxOffset,
		value:    0,
	}
}

func (c Cursor) Value() int { return c.value }
func (c *Cursor) Set(v int) error {
	if err := c.checkOffset(v); err != nil {
		return fmt.Errorf("new offset value is invalid: %w", err)
	}

	c.value = v
	return nil
}

func (c Cursor) checkOffset(v int) error {
	if v < 0 {
		return fmt.Errorf("offset cannot be negative: %v", v)
	}
	if v >= c.maxValue {
		return fmt.Errorf("offset is too high: %v >= %v", v, c.maxValue)
	}

	return nil
}
