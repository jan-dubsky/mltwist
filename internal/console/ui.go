package console

import (
	"decomp/internal/console/internal/control"
	"decomp/internal/console/internal/cursor"
	"decomp/internal/console/internal/lines"
	"decomp/internal/console/internal/view"
	"decomp/internal/deps"
	"errors"
	"fmt"
)

type UI struct {
	view    printer
	control controller
}

func NewUI(p *deps.Program) *UI {
	lines := lines.NewLines(p)
	cursor := cursor.New(lines)

	view := view.New(lines, cursor)
	return &UI{
		view:    view,
		control: control.New(lines, cursor, view),
	}
}

func (ui *UI) Run() error {
	for {
		err := ui.view.Print()
		if err != nil {
			return fmt.Errorf("cannot print UI output: %w", err)
		}

		fmt.Printf("\n")

		err = ui.control.Command()
		if err != nil {
			if errors.Is(err, control.ErrQuit) {
				return nil
			}

			return fmt.Errorf("cannot process command: %w", err)
		}
	}
}
