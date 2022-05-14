package disassemble

import (
	"fmt"
	"math"
	"mltwist/internal/consoleui"
	"mltwist/internal/consoleui/internal/cmdtools"
	"mltwist/internal/consoleui/internal/linereader"
	"mltwist/internal/consoleui/internal/lines"
	"regexp"
)

func commands(m *mode) []consoleui.Command {
	return []consoleui.Command{{
		Keys: []string{"down", "d"},
		Help: "Move line cursor <N> lines down.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			return m.view.Cursor.Set(m.view.Cursor.Value() + args[0].(int))
		},
	}, {
		Keys: []string{"up", "u"},
		Help: "Move line cursor <N> lines up.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			return m.view.Cursor.Set(m.view.Cursor.Value() + -args[0].(int))
		},
	}, {
		Keys: []string{"move", "mv", "m"},
		Help: "Move instruction from line <N> to line <M>",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			from, to := args[0].(int), args[1].(int)
			m.view.Lines.UnmarkAll()

			err := m.view.Lines.Move(from, to)
			if err != nil {
				m.view.Lines.SetMark(from, lines.MarkErrMovedFrom)
				m.view.Lines.SetMark(to, lines.MarkErrMovedTo)
				return err
			}

			m.view.Lines.SetMark(from, lines.MarkMovedFrom)
			m.view.Lines.SetMark(to, lines.MarkMovedTo)
			return nil
		},
	}, {
		Keys: []string{"bounds", "b"},
		Help: "Show bounds where a given instruction can be moved.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			l := args[0].(int)
			m.view.Lines.UnmarkAll()

			block, ok := m.view.Lines.Block(l)
			if !ok {
				m.view.Lines.SetMark(l, lines.MarkErr)
				return fmt.Errorf("line doesn't belong to a block: %d", l)
			}

			ins, ok := m.view.Lines.Index(l).Instruction()
			if !ok {
				m.view.Lines.SetMark(l, lines.MarkErr)
				return fmt.Errorf("line is not an instruction: %d", l)
			}

			lower := block.LowerBound(ins)
			upper := block.UpperBound(ins)
			lowerLine := m.view.Lines.Line(block, lower)
			upperLine := m.view.Lines.Line(block, upper)

			// Lower and Upper indices are inclusive, but in
			// visualization we want to have exclusive indices.
			m.view.Lines.SetMark(lowerLine-1, lines.MarkLowerBound)
			m.view.Lines.SetMark(upperLine+1, lines.MarkUpperBound)

			return nil
		},
	}, {
		Keys: []string{"find", "f", "/"},
		Help: "Find row matching standard POSIX regex.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseString,
		},
		OptionalArgs: cmdtools.JoinOptStrings,
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			r := args[0].(string)
			if len(args) > 1 {
				r = r + " " + args[1].(string)
			}

			regexp, err := regexp.CompilePOSIX(r)
			if err != nil {
				return fmt.Errorf("invalid regex %q: %w", r, err)
			}

			var line int = -1
			offset := m.view.Cursor.Value()
			for i := offset + 1; i != offset; i = (i + 1) % m.view.Lines.Len() {
				if regexp.MatchString(m.view.Lines.Index(i).String()) {
					line = i
					break
				}
			}

			if line == -1 {
				return linereader.ErrMsgf(
					"No line matching regex %q found.\n", r)
			}

			err = m.view.Cursor.Set(line)
			if err != nil {
				return err
			}
			return nil
		},
	}, {
		Keys: []string{"goto", "g"},
		Help: "Go to line number <N>.",
		Args: []consoleui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			n := args[0].(int)
			if l := m.view.Lines.Len(); n > l {
				return fmt.Errorf("line number too big: %d > %d", n, l)
			}

			if err := m.view.Cursor.Set(n); err != nil {
				return err
			}

			return nil
		},
	}, {
		Keys: []string{"entrypoint", "entry"},
		Help: "Sets cursor to app entrypoint.",
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			a := m.prog.Entrypoint()
			block, ok := m.prog.Address(a)
			if !ok {
				return fmt.Errorf("cannot find block at address 0x%x", a)
			}

			ins, ok := block.Address(a)
			if !ok {
				return fmt.Errorf("cannot find instruction at address 0x%x", a)
			}

			line := m.view.Lines.Line(block, ins.Idx())
			m.view.Cursor.Set(line)
			return nil
		},
	}, {
		Keys: []string{"alllines"},
		Help: "Prints all lines of the code into console. " +
			"Ignores current cursor position.",
		Action: func(_ *consoleui.UI, args ...interface{}) error {
			for i := 0; i < m.view.Lines.Len(); i++ {
				fmt.Print(m.view.Format(i))
			}

			if err := linereader.ErrMsgf("\n"); err != nil {
				return err
			}

			return nil
		},
	}, {
		Keys: []string{"emulate", "emul", "e"},
		Help: "Start emulating the machine code at current line.\n\n" +
			"TIP: for emulation started at entrypoint, use 'entrypoint' " +
			"command followed by this command.",
		Action: func(ui *consoleui.UI, args ...interface{}) error {
			l := m.view.Cursor.Value()
			line := m.view.Lines.Index(l)

			block, ok := m.view.Lines.Block(l)
			if !ok {
				return fmt.Errorf("line %d belongs to no block", l)
			}

			insIdx, ok := line.Instruction()
			if !ok {
				return fmt.Errorf("line %d is not instruction line", l)
			}

			ins := block.Index(insIdx)
			emul, err := m.emulFunc(m.prog, ins.Addr())
			if err != nil {
				return fmt.Errorf("bug: cannot create emulation: %w", err)
			}

			return ui.AddMode("emulate", emul)
		},
	},
	}
}
