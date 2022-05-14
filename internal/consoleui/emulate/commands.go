package emulate

import (
	"fmt"
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/cmdtools"
	"mltwist/internal/consoleui/internal/linereader"
	"mltwist/internal/consoleui/internal/memview"
	"mltwist/pkg/expr"
	"sort"
)

func commands(m *mode) []consoleui.Command {
	return []consoleui.Command{{
		Keys: []string{"forward", "fwd", "f", "step", "s"},
		Help: "Move emulation one instruction forward.",
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			_, err := m.emul.Step()
			if err != nil {
				return fmt.Errorf("cannot emulate instruction: %w", err)
			}

			err = m.refreshCursor()
			if err != nil {
				return fmt.Errorf("cannot refresh cursor: %w", err)
			}

			return nil
		},
	}, {
		Keys: []string{"memories", "mems", "ms"},
		Help: "List all memories the program wrote.",
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			mems := m.stat.Mems
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
		Args: []consoleui.ArgParseFunc{cmdtools.ParseString},
		Action: func(ui *consoleui.UI, args ...interface{}) error {
			key := expr.Key(args[0].(string))
			mem := m.stat.Mems[key]

			mode := memview.New(mem)
			name := fmt.Sprintf("memview(%s)", key)
			return ui.AddMode(name, mode)
		},
	}, {
		Keys: []string{"reqmod", "rmod"},
		Help: "Modify register of the running application.",
		Args: []consoleui.ArgParseFunc{cmdtools.ParseString},
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
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			key := expr.Key(args[0].(string))
			regs := m.stat.Regs

			r, ok := regs[key]
			if !ok {
				return fmt.Errorf("register is not set at all: %s", key)
			}

			regs[key] = readRegister(key, r.Width())
			return nil
		},
	}}
}
