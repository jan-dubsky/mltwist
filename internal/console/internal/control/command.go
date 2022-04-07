package control

import (
	"fmt"
	"strconv"
	"strings"
)

type argParseFunc func(s string) (interface{}, error)
type optArgParseFunc func(s []string) ([]interface{}, error)

type command struct {
	keys         []string
	help         string
	args         []argParseFunc
	optionalArgs optArgParseFunc
	action       func(c *Control, args ...interface{}) error
}

func (c command) keysString() string { return strings.Join(c.keys, ", ") }

// parseNum parses an integer parameter out of string and then validates that
// the value is in between min and max. Both min and max are inclusive
// boundaries.
func parseNum(min, max int) argParseFunc {
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

func parseString(s string) (interface{}, error) {
	return s, nil
}

func joinOptStrings(strs []string) ([]interface{}, error) {
	return []interface{}{strings.Join(strs, "")}, nil
}
