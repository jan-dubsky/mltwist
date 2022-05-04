package emulate

import (
	"fmt"
	"mltwist/internal/console/ui"
)

func commands(e *emulate) []ui.Command {
	return []ui.Command{{
		Keys: []string{"forward", "fwd", "f", "step", "s"},
		Help: "Move emulation one instruction forward.",
		Action: func(c *ui.Control, args ...interface{}) error {
			ip := e.emul.IP()
			ins, ok := e.p.AddrIns(ip)
			if !ok {
				return fmt.Errorf(
					"cannot find instruction at address 0x%x", ip)
			}

			_ = e.emul.Step(ins)

			err := e.refreshCursor()
			if err != nil {
				return fmt.Errorf("cannot refresh cursor: %w", err)
			}

			return nil
		},
	}}
}
