package disassemble

import (
	"mltwist/internal/console/internal/cursor"
	"mltwist/internal/console/internal/lines"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/console/ui"
	"mltwist/internal/deps"
)

var _ ui.Mode = &mode{}

type mode struct {
	prog   *deps.Program
	lines  *lines.Lines
	cursor *cursor.Cursor
	view   *view.LinesView
}

func New(p *deps.Program) *mode {
	lines := lines.New(p)
	cursor := cursor.New(lines.Len())

	return &mode{
		prog:   p,
		lines:  lines,
		cursor: cursor,
		view:   view.NewLinesView(lines, cursor),
	}
}

func (d *mode) Commands() []ui.Command { return commands(d) }
func (d *mode) View() view.Element     { return d.view }
