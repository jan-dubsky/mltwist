package memview

import (
	"math"
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/cmdtools"
)

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
	}}
}
