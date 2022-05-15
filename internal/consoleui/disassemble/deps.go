package disassemble

import (
	"mltwist/internal/consoleui"
	"mltwist/internal/deps"
	"mltwist/pkg/model"
)

// EmulFunc is a function creating program emulation from current value of a
// program and current value of an instruction pointer.
type EmulFunc func(p *deps.Code, ip model.Addr) (consoleui.Mode, error)
