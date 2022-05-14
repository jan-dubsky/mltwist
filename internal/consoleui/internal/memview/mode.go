package memview

import (
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/view"
	"mltwist/internal/state/memory"
)

var _ consoleui.Mode = &mode{}

type mode struct {
	view *memoryView
}

func New(mem memory.Memory) *mode {
	return &mode{
		view: newMemoryView(mem),
	}
}

func (m *mode) Commands() []consoleui.Command { return commands(m) }
func (m *mode) View() view.View               { return m.view }
