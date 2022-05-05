package memview

import (
	"mltwist/internal/console/internal/view"
	"mltwist/internal/console/ui"
	"mltwist/internal/state"
)

var _ ui.Mode = &mode{}

type mode struct {
	view *memoryView
}

func New(mem *state.Memory) *mode {
	return &mode{
		view: newMemoryView(mem),
	}
}

func (m *mode) Commands() []ui.Command { return commands(m) }

func (m *mode) View() view.Element { return m.view }
