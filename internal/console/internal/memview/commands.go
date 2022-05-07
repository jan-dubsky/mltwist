package memview

import (
	"math"
	"mltwist/internal/console/internal/ui"
	"mltwist/internal/console/internal/ui/cmdtools"
)

func commands(m *mode) []ui.Command {
	return []ui.Command{{
		Keys: []string{"down", "d"},
		Help: "Move line cursor <N> lines down.",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *ui.UI, args ...interface{}) error {
			return m.view.c.Set(m.view.c.Value() + args[0].(int))
		},
	}, {
		Keys: []string{"up", "u"},
		Help: "Move line cursor <N> lines up.",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *ui.UI, args ...interface{}) error {
			return m.view.c.Set(m.view.c.Value() - args[0].(int))
		},
	}}
}
