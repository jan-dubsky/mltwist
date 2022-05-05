package emulate

import (
	"fmt"
	"mltwist/internal/console/emulate/memview"
	"mltwist/internal/console/internal/linereader"
	"mltwist/internal/console/ui"
	"mltwist/internal/console/ui/cmdtools"
	"mltwist/pkg/expr"
	"sort"
)

func commands(m *mode) []ui.Command {
	return []ui.Command{{
		Keys: []string{"forward", "fwd", "f", "step", "s"},
		Help: "Move emulation one instruction forward.",
		Action: func(c *ui.Control, args ...interface{}) error {
			ip := m.emul.IP()
			ins, ok := m.p.AddrIns(ip)
			if !ok {
				return fmt.Errorf(
					"cannot find instruction at address 0x%x", ip)
			}

			_ = m.emul.Step(ins)

			err := m.refreshCursor()
			if err != nil {
				return fmt.Errorf("cannot refresh cursor: %w", err)
			}

			return nil
		},
	}, {
		Keys: []string{"memories", "mems", "ms"},
		Help: "List all memories the program wrote.",
		Action: func(c *ui.Control, args ...interface{}) error {
			mems := m.emul.State().Mems
			keys := make([]expr.Key, 0, len(mems))
			for k := range mems {
				keys = append(keys, k)
			}

			sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

			fmt.Printf("program wrote %d memories:\n", len(keys))
			for _, k := range keys {
				fmt.Printf("\t%s\n", k)
			}

			fmt.Printf("\nPress ENTER to continue\n")
			if _, err := linereader.ReadLine(); err != nil {
				return fmt.Errorf("readline error: %w", err)
			}

			return nil
		},
	}, {
		Keys: []string{"memory", "mem", "m"},
		Help: "Show content of memory address space identified by <key>.",
		Args: []ui.ArgParseFunc{cmdtools.ParseString},
		Action: func(c *ui.Control, args ...interface{}) error {
			key := expr.Key(args[0].(string))
			mem := m.emul.State().Mems[key]

			mode := memview.New(mem)
			name := fmt.Sprintf("memview(%s)", key)
			return c.AddMode(name, mode)
		},
	}, {
		Keys: []string{"reqmod", "rmod"},
		Help: "Modify register of the running application.",
		Args: []ui.ArgParseFunc{cmdtools.ParseString},
		// Please note that it makes absolutely no sense to change the
		// register width. The register width is given by instructions
		// preceding the current instruction. Consequently the register
		// width is an inherent property of the program -> it cannot be
		// changed from outside.
		//
		// Please note that the statement above about inherent property
		// of the program is not right. Due to indirect jump
		// instructions, values on inputs of the program are allowed to
		// effect order of instruction (basic block) execution in the
		// program. Consequently the register width is inherent property
		// of this specific program run. This fact doesn't weaken the
		// final statement that register width cannot be changed
		// dynamically.
		Action: func(c *ui.Control, args ...interface{}) error {
			key := expr.Key(args[0].(string))
			regs := m.emul.State().Regs

			r, ok := regs[key]
			if !ok {
				return fmt.Errorf("register is not set at all: %s", key)
			}

			regs[key] = readRegister(key, r.Width())
			return nil
		},
	}}
}
