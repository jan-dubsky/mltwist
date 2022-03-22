package console

import (
	"decomp/internal/console/internal/control"
	"decomp/internal/console/internal/lines"
	"decomp/internal/console/internal/view"
	"decomp/internal/deps"
	"fmt"
)

type UI struct {
	view    printer
	control controller
}

func NewUI(p *deps.Program) *UI {
	lines := lines.NewLines(p)
	return &UI{
		view:    view.New(lines),
		control: control.New(lines),
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
			return fmt.Errorf("cannot process command: %w", err)
		}
	}
}
