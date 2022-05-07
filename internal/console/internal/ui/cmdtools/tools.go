package cmdtools

import (
	"fmt"
	"mltwist/internal/console/internal/ui"
	"strconv"
	"strings"
)

// ParseNum parses an integer parameter out of string and then validates that
// the value is in between min and max. Both min and max are inclusive
// boundaries.
func ParseNum(min, max int) ui.ArgParseFunc {
	return func(s string) (interface{}, error) {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value %q: %w", s, err)
		}

		if v < min {
			return nil, fmt.Errorf(
				"value is less than allowed minimum: %d < %d", v, min)
		}
		if v > max {
			return nil, fmt.Errorf(
				"value is greater than allowed maximum: %d > %d", v, max)
		}

		return v, nil
	}
}

func ParseString(s string) (interface{}, error) {
	return s, nil
}

func JoinOptStrings(strs []string) ([]interface{}, error) {
	return []interface{}{strings.Join(strs, "")}, nil
}
