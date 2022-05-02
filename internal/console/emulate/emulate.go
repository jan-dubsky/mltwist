package emulate

import (
	"fmt"
	"mltwist/internal/console/internal/cursor"
	"mltwist/internal/console/internal/lines"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/console/ui"
	"mltwist/internal/deps"
	"mltwist/internal/emulator"
	"mltwist/pkg/model"
)

var _ ui.Mode = &emulate{}

type emulate struct {
	p    *deps.Program
	emul *emulator.Emulator

	lines  *lines.Lines
	cursor *cursor.Cursor

	view view.Element
}

func New(p *deps.Program, ip model.Addr) (*emulate, error) {
	emul := emulator.New(p, ip, &stateProvider{})

	lines := lines.New(p)
	cursor := cursor.New(lines)

	lineView := view.NewLinesView(lines, cursor)
	regView := newValuesElement(emul.State())

	e := &emulate{
		p:    p,
		emul: emul,

		lines:  lines,
		cursor: cursor,

		view: view.NewCompositeView(lineView, regView),
	}

	if err := e.refreshCursor(); err != nil {
		return nil, fmt.Errorf("cannot set cursor: %w", err)
	}

	return e, nil
}

func (e *emulate) Commands() []ui.Command { return commands(e) }
func (e *emulate) Element() view.Element  { return e.view }

func (e *emulate) refreshCursor() error {
	ip := e.emul.IP()

	block, ok := e.p.Addr(ip)
	if !ok {
		return fmt.Errorf("cannot find block containing address 0x%x", ip)
	}

	ins, ok := block.Addr(ip)
	if !ok {
		return fmt.Errorf("cannot find instruction at address 0x%x", ip)
	}

	e.cursor.Set(e.lines.Line(block, ins.Idx()))
	return nil
}
