package disassemble

import (
	"mltwist/internal/console/internal/cursor"
	"mltwist/internal/console/internal/lines"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/console/ui"
	"mltwist/internal/deps"
)

var _ ui.Mode = &disassemble{}

type disassemble struct {
	prog   *deps.Program
	lines  *lines.Lines
	cursor *cursor.Cursor
	view   *view.LinesView
}

func New(p *deps.Program) *disassemble {
	lines := lines.New(p)
	cursor := cursor.New(lines)

	return &disassemble{
		prog:   p,
		lines:  lines,
		cursor: cursor,
		view:   view.NewLinesView(lines, cursor),
	}
}

func (d *disassemble) Commands() []ui.Command { return commands(d) }
func (d *disassemble) Element() view.Element  { return d.view }
