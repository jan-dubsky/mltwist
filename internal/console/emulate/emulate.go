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

var _ ui.Mode = &mode{}

type mode struct {
	p    *deps.Program
	emul *emulator.Emulator

	lines  *lines.Lines
	cursor *cursor.Cursor

	view view.Element
}

func New(p *deps.Program, ip model.Addr) (*mode, error) {
	emul := emulator.New(p, ip, &stateProvider{})

	lines := lines.New(p)
	cursor := cursor.New(lines.Len())

	lineView := view.NewLinesView(lines, cursor)
	regView := newRegView(emul.State())

	e := &mode{
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

func (e *mode) Commands() []ui.Command { return commands(e) }
func (e *mode) View() view.Element     { return e.view }

func (e *mode) refreshCursor() error {
	ip := e.emul.IP()

	block, ok := e.p.Address(ip)
	if !ok {
		return fmt.Errorf("cannot find block containing address 0x%x", ip)
	}

	ins, ok := block.Address(ip)
	if !ok {
		return fmt.Errorf("cannot find instruction at address 0x%x", ip)
	}

	e.cursor.Set(e.lines.Line(block, ins.Idx()))
	return nil
}
