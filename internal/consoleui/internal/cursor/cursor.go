package cursor

import (
	"fmt"
)

// Cursor is a simple gadget struct storing index in range [0, maxValue).
type Cursor struct {
	maxValue int
	value    int
}

// New creates a new cursor for values [0, maxValue).
func New(maxValue int) *Cursor {
	return &Cursor{
		maxValue: maxValue,
		value:    0,
	}
}

// Value returns the current value stored in cursor.
func (c Cursor) Value() int { return c.value }

// MaxValue returns maximal exclusive value which can be stored in the cursor.
func (c Cursor) MaxValue() int { return c.maxValue }

// Set changes value of cursor to v. If v is less than zero or grater or equal
// to MaxValue, this function will not change the cursor value. Instead, it will
// return an error.
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
