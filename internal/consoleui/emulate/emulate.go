package emulate

import (
	"fmt"
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/lines"
	"mltwist/internal/consoleui/internal/view"
	"mltwist/internal/deps"
	"mltwist/internal/emulator"
	"mltwist/internal/state"
	"mltwist/pkg/model"
)

var _ consoleui.Mode = &mode{}

type mode struct {
	code *deps.Code
	stat *state.State
	emul *emulator.Emulator

	lineView *lines.View
	view     view.View
}

func New(code *deps.Code, ip model.Addr, stat *state.State) (*mode, error) {
	emul := emulator.New(code, ip, stat, &stateProvider{})

	lineView := lines.NewView(code)
	regView := newRegView(stat)

	e := &mode{
		code: code,
		stat: stat,
		emul: emul,

		lineView: lineView,
		view:     view.NewCompositeView(lineView, regView),
	}

	if err := e.refreshCursor(); err != nil {
		return nil, fmt.Errorf("cannot set cursor: %w", err)
	}

	return e, nil
}

func (e *mode) Commands() []consoleui.Command { return commands(e) }
func (e *mode) View() view.View               { return e.view }

func (e *mode) refreshCursor() error {
	ip := e.emul.IP()

	block, ok := e.code.Address(ip)
	if !ok {
		return fmt.Errorf("cannot find block containing address 0x%x", ip)
	}

	ins, ok := block.Address(ip)
	if !ok {
		return fmt.Errorf("cannot find instruction at address 0x%x", ip)
	}

	e.lineView.Cursor.Set(e.lineView.Lines.Line(block, ins.Idx()))
	return nil
}
