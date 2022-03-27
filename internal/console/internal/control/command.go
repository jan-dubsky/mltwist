package control

import (
	"decomp/internal/console/internal/lines"
	"decomp/internal/deps"
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

// sequentialOptArgParseFunc is a helper optional argument parses which uses f
// to parse every individual optional argument and returns an array with all
// parsed values.
func sequentialOptArgParseFunc(f argParseFunc) optArgParseFunc {
	return func(strs []string) ([]interface{}, error) {
		vals := make([]interface{}, len(strs))
		for i, s := range strs {
			v, err := f(s)
			if err != nil {
				return nil, fmt.Errorf("optional arg %d/%d: %w",
					i, len(strs), err)
			}

			vals[i] = v
		}

		return vals, nil
	}
}

func joinOptStrings(strs []string) ([]interface{}, error) {
	return []interface{}{strings.Join(strs, "")}, nil
}

// insLine resurns basic-block and an instruction the line is referring. This
// method failed with an error if line lineIdx in lines is not an instruction
// line.
func insLine(lines *lines.Lines, lineIdx int) (deps.Block, int, error) {
	block, err := blockLine(lines, lineIdx)
	if err != nil {
		return deps.Block{}, -1, err
	}

	ins, ok := lines.Index(lineIdx).Instruction()
	if !ok {
		err := fmt.Errorf("line %d is not an instruction", lineIdx)
		return deps.Block{}, -1, err
	}

	return block, ins, nil
}

func blockLine(lines *lines.Lines, lineIdx int) (deps.Block, error) {
	block, ok := lines.Block(lineIdx)
	if !ok {
		err := fmt.Errorf("line %d doesn't belong to any block", lineIdx)
		return deps.Block{}, err
	}

	return block, nil
}
