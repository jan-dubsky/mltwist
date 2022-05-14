package disassemble

import (
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/lines"
	"mltwist/internal/consoleui/internal/view"
	"mltwist/internal/deps"
)

var _ consoleui.Mode = &mode{}

type mode struct {
	prog *deps.Program
	view *lines.View

	emulFunc EmulFunc
}

// New creates a new disassembler UI mode displaying and manipulating
// instructions from p.
func New(p *deps.Program, emulF EmulFunc) consoleui.Mode {
	return &mode{
		prog:     p,
		view:     lines.NewView(p),
		emulFunc: emulF,
	}
}

func (d *mode) Commands() []consoleui.Command { return commands(d) }
func (d *mode) View() view.View               { return d.view }
