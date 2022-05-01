package disassemble

import (
	"fmt"
	"math"
	"mltwist/internal/console/internal/lines"
	"mltwist/internal/console/ui"
	"mltwist/internal/console/ui/cmdtools"
	"regexp"
)

func commands(d *disassemble) []ui.Command {
	return []ui.Command{{
		Keys: []string{"down", "d"},
		Help: "Move line cursor <N> lines down.",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(c *ui.Control, args ...interface{}) error {
			return d.c.Set(d.c.Value() + args[0].(int))
		},
	}, {
		Keys: []string{"up", "u"},
		Help: "Move line cursor <N> lines up.",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(c *ui.Control, args ...interface{}) error {
			return d.c.Set(d.c.Value() + -args[0].(int))
		},
	}, {
		Keys: []string{"move", "mv", "m"},
		Help: "Move instruction from line <N> to line <M>",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(c *ui.Control, args ...interface{}) error {
			from, to := args[0].(int), args[1].(int)
			d.l.UnmarkAll()

			err := d.l.Move(from, to)
			if err != nil {
				d.l.SetMark(from, lines.MarkErrMovedFrom)
				d.l.SetMark(to, lines.MarkErrMovedTo)
				return err
			}

			d.l.SetMark(from, lines.MarkMovedFrom)
			d.l.SetMark(to, lines.MarkMovedTo)
			return nil
		},
	}, {
		Keys: []string{"bounds", "b"},
		Help: "Show bounds where a given instruction can be moved.",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(c *ui.Control, args ...interface{}) error {
			l := args[0].(int)
			d.l.UnmarkAll()

			block, ok := d.l.Block(l)
			if !ok {
				d.l.SetMark(l, lines.MarkErr)
				return fmt.Errorf("line doesn't belong to a block: %d", l)
			}

			ins, ok := d.l.Index(l).Instruction()
			if !ok {
				d.l.SetMark(l, lines.MarkErr)
				return fmt.Errorf("line is not an instruction: %d", l)
			}

			lower := block.LowerBound(ins)
			upper := block.UpperBound(ins)
			lowerLine := d.l.Line(block, lower)
			upperLine := d.l.Line(block, upper)

			// Lower and Upper indices are inclusive, but in
			// visualization we want to have exclusive indices.
			d.l.SetMark(lowerLine-1, lines.MarkLowerBound)
			d.l.SetMark(upperLine+1, lines.MarkUpperBound)

			return nil
		},
	}, {
		Keys: []string{"find", "f", "/"},
		Help: "Find row matching standard POSIX regex.",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseString,
		},
		OptionalArgs: cmdtools.JoinOptStrings,
		Action: func(c *ui.Control, args ...interface{}) error {
			r := args[0].(string)
			if len(args) > 1 {
				r = r + " " + args[1].(string)
			}

			regexp, err := regexp.CompilePOSIX(r)
			if err != nil {
				return fmt.Errorf("invalid regex %q: %w", r, err)
			}

			var line int = -1
			offset := d.c.Value()
			for i := offset + 1; i != offset; i = (i + 1) % d.l.Len() {
				if regexp.MatchString(d.l.Index(i).String()) {
					line = i
					break
				}
			}

			if line == -1 {
				return c.ErrMsgf("No line matching regex %q found.\n", r)
			}

			err = d.c.Set(line)
			if err != nil {
				return err
			}
			return nil
		},
	}, {
		Keys: []string{"goto", "g"},
		Help: "Go to line number <N>.",
		Args: []ui.ArgParseFunc{
			cmdtools.ParseNum(0, math.MaxInt),
		},
		Action: func(c *ui.Control, args ...interface{}) error {
			n := args[0].(int)
			if l := d.l.Len(); n > l {
				return fmt.Errorf("line number too big: %d > %d", n, l)
			}

			if err := d.c.Set(n); err != nil {
				return err
			}

			return nil
		},
	}, {
		Keys: []string{"entrypoint", "entry"},
		Help: "Sets cursor to app entrypoint.",
		Action: func(c *ui.Control, args ...interface{}) error {
			a := d.p.Entrypoint()
			block, ok := d.p.Addr(a)
			if !ok {
				return fmt.Errorf("cannot find block at address 0x%x", a)
			}

			ins, ok := block.Addr(a)
			if !ok {
				return fmt.Errorf("cannot find instruction at address 0x%x", a)
			}

			line := d.l.Line(block, ins.Idx())
			d.c.Set(line)
			return nil
		},
	}, {
		Keys: []string{"alllines"},
		Help: "Prints all lines of the code into console. " +
			"Ignores current cursor position.",
		Action: func(c *ui.Control, args ...interface{}) error {
			for i := 0; i < d.l.Len(); i++ {
				fmt.Print(d.v.Format(i))
			}

			if err := c.ErrMsgf("\n"); err != nil {
				return err
			}

			return nil
		},
	}, /*		{
			Keys: []string{"emulate", "emul", "e"},
			Help: "Start emulating the machine code at current line.",
			Action: func(c *control.Control, args ...interface{}) error {
				return c.AddMode("emulate", listEmulate())
			},
		},*/
	}
}
