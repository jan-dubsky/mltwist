package disassemble

import (
	"mltwist/internal/console/internal/lines"
	"mltwist/internal/console/internal/ui"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/deps"
)

var _ ui.Mode = &mode{}

type mode struct {
	prog *deps.Program
	view *lines.View
}

// New creates a new disassembler UI mode displaying and manipulating
// instructions from p.
func New(p *deps.Program) ui.Mode {
	return &mode{
		prog: p,
		view: lines.NewView(p),
	}
}

func (d *mode) Commands() []ui.Command { return commands(d) }
func (d *mode) View() view.View        { return d.view }
