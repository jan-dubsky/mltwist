package disassemble

import (
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/lines"
	"mltwist/internal/consoleui/internal/view"
	"mltwist/internal/deps"
)

var _ consoleui.Mode = &mode{}

type mode struct {
	code *deps.Code
	view *lines.View

	emulFunc EmulFunc
}

// New creates a new disassembler UI mode displaying and manipulating
// instructions from p.
func New(code *deps.Code, emulF EmulFunc) consoleui.Mode {
	return &mode{
		code:     code,
		view:     lines.NewView(code),
		emulFunc: emulF,
	}
}

func (d *mode) Commands() []consoleui.Command { return commands(d) }
func (d *mode) View() view.View               { return d.view }
