package console

import (
	"fmt"
	"mltwist/internal/console/internal/disassemble"
	"mltwist/internal/console/internal/ui"
	"mltwist/internal/deps"
)

type Console struct {
	ui *ui.UI
}

func New(p *deps.Program) (*Console, error) {
	ui, err := ui.New(disassemble.New(p))
	if err != nil {
		return nil, fmt.Errorf("cannot create UI: %w", err)
	}

	return &Console{
		ui: ui,
	}, nil
}

func (c *Console) Run() error { return c.ui.Run() }
