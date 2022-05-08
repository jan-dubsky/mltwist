package memview

import (
	"mltwist/internal/console/internal/ui"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/state/memory"
)

var _ ui.Mode = &mode{}

type mode struct {
	view *memoryView
}

func New(mem memory.Memory) *mode {
	return &mode{
		view: newMemoryView(mem),
	}
}

func (m *mode) Commands() []ui.Command { return commands(m) }

func (m *mode) View() view.View { return m.view }
