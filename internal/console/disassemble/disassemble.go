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
	p *deps.Program
	l *lines.Lines
	c *cursor.Cursor
	v *view.LinesView
}

func New(p *deps.Program) *disassemble {
	lines := lines.New(p)
	cursor := cursor.New(lines)

	return &disassemble{
		p: p,
		l: lines,
		c: cursor,
		v: view.NewLinesView(lines, cursor),
	}
}

func (d *disassemble) Commands() []ui.Command { return commands(d) }
func (d *disassemble) Element() view.Element  { return d.v }
