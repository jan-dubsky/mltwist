package console

import (
	"errors"
	"fmt"
	"mltwist/internal/console/internal/control"
	"mltwist/internal/console/internal/cursor"
	"mltwist/internal/console/internal/lines"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/deps"
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
		control: control.New(p, lines, cursor, view),
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
