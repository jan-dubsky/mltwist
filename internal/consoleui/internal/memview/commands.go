package memview

import (
	"fmt"
	"math"
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/cmdtools"
	"mltwist/pkg/model"
	"strconv"
	"unsafe"
)

func parseAddr(s string) (interface{}, error) {
	base := 10
	if len(s) > 2 && s[:2] == "0x" || s[:2] == "0X" {
		base = 16
		s = s[2:]
	} else if len(s) == 2 && s[:2] == "0b" || s[:2] == "0B" {
		base = 2
		s = s[2:]
	} else if len(s) > 0 && s[0] == '0' {
		base = 8
		s = s[1:]
	}

	addr, err := strconv.ParseUint(s, base, int(unsafe.Sizeof(model.Addr(0))*8))
	if err != nil {
		return nil, fmt.Errorf("invalid uint with base %d: %w", base, err)
	}

	return model.Addr(addr), nil
}

func commands(m *mode) []consoleui.Command {
	return []consoleui.Command{{
		Keys: []string{"down", "d"},
		Help: "Move line cursor <N> lines down.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			return m.view.c.Set(m.view.c.Value() + args[0].(int))
		},
	}, {
		Keys: []string{"up", "u"},
		Help: "Move line cursor <N> lines up.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			return m.view.c.Set(m.view.c.Value() - args[0].(int))
		},
	}, {
		Keys: []string{"goto", "g"},
		Help: "Go to line <n>.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			line := args[0].(int)
			return m.view.c.Set(line)
		},
	}, {
		Keys: []string{"address", "addr", "a"},
		Help: "Go to address.",
		Args: []consoleui.ArgParseFunc{parseAddr},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			addr := args[0].(model.Addr)

			idx := -1
		iterLines:
			for i, l := range m.view.lines {
				for _, r := range l.ranges {
					if r.Containts(addr) {
						idx = i
						break iterLines
					}
				}
			}

			if idx < 0 {
				return fmt.Errorf("no line with address 0x%x found", addr)
			}

			err := m.view.c.Set(idx)
			if err != nil {
				return fmt.Errorf("cannot set cursor to %d: %w", idx, err)
			}

			return nil
		},
	}}
}
