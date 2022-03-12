package control

import (
	"decomp/internal/console/internal/lines"
	"decomp/internal/deps"
	"fmt"
	"strconv"
	"strings"
)

type argParseFunc func(s string) (interface{}, error)

type command struct {
	keys         []string
	help         string
	args         []argParseFunc
	optionalArgs []argParseFunc
	action       func(c *Control, args ...interface{}) error
}

func (c command) keysString() string { return strings.Join(c.keys, ", ") }

// parseNum parses an integer parameter out of string and then validates that
// the value is in between min and max. Both min and max are inclusive
// boundaries.
func parseNum(min, max int) argParseFunc {
	return func(s string) (interface{}, error) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value %q: %w", s, err)
		}

		if val < min {
			return nil, fmt.Errorf("value is less than allowed minimum: %d < %d", val, min)
		}
		if val > max {
			return nil, fmt.Errorf("value is greater than allowed maximum: %d > %d", val, max)
		}

		return val, nil
	}
}

func insLine(lines *lines.Lines, lineIdx int) (*deps.Block, deps.Instruction, error) {
	block, ok := lines.Block(lineIdx)
	if !ok {
		err := fmt.Errorf("line %d doesn't belong to any block", lineIdx)
		return nil, deps.Instruction{}, err
	}

	ins, ok := lines.Instruction(lineIdx)
	if !ok {
		err := fmt.Errorf("line %d is not an instruction", lineIdx)
		return nil, deps.Instruction{}, err
	}

	return block, ins, nil
}
