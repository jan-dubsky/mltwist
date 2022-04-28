package control

import (
	"fmt"
	"math"
	"mltwist/internal/console/internal/lines"
	"regexp"
)

// ErrQuit is returned by command in case UI termination is required.
var ErrQuit = fmt.Errorf("app exit required")

var standardCmds = []*command{{
	keys: []string{"quit", "q"},
	help: "Quit the app.",
	action: func(c *Control, args ...interface{}) error {
		fmt.Printf("\n")
		return ErrQuit
	}},
}

func listDisassemble() []*command {
	return []*command{{
		keys: []string{"down", "d"},
		help: "Move line cursor <N> lines down.",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			return c.c.Set(c.c.Value() + args[0].(int))
		},
	}, {
		keys: []string{"up", "u"},
		help: "Move line cursor <N> lines up.",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			return c.c.Set(c.c.Value() + -args[0].(int))
		},
	}, {
		keys: []string{"move", "mv", "m"},
		help: "Move instruction from line <N> to line <M>",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			from, to := args[0].(int), args[1].(int)
			c.l.UnmarkAll()

			err := c.l.Move(from, to)
			if err != nil {
				c.l.SetMark(from, lines.MarkErrMovedFrom)
				c.l.SetMark(to, lines.MarkErrMovedTo)
				return err
			}

			c.l.SetMark(from, lines.MarkMovedFrom)
			c.l.SetMark(to, lines.MarkMovedTo)
			return nil
		},
	}, {
		keys: []string{"bounds", "b"},
		help: "Show bounds where a given instruction can be moved.",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			l := args[0].(int)
			c.l.UnmarkAll()

			block, ok := c.l.Block(l)
			if !ok {
				c.l.SetMark(l, lines.MarkErr)
				return fmt.Errorf("line doesn't belong to a block: %d", l)
			}

			ins, ok := c.l.Index(l).Instruction()
			if !ok {
				c.l.SetMark(l, lines.MarkErr)
				return fmt.Errorf("line is not an instruction: %d", l)
			}

			lower := block.LowerBound(ins)
			upper := block.UpperBound(ins)
			lowerLine := c.l.Line(block, lower)
			upperLine := c.l.Line(block, upper)

			// Lower and Upper indices are inclusive, but in
			// visualization we want to have exclusive indices.
			c.l.SetMark(lowerLine-1, lines.MarkLowerBound)
			c.l.SetMark(upperLine+1, lines.MarkUpperBound)

			return nil
		},
	}, {
		keys: []string{"find", "f", "/"},
		help: "Find row matching standard POSIX regex.",
		args: []argParseFunc{
			parseString,
		},
		optionalArgs: joinOptStrings,
		action: func(c *Control, args ...interface{}) error {
			r := args[0].(string)
			if len(args) > 1 {
				r = r + " " + args[1].(string)
			}

			regexp, err := regexp.CompilePOSIX(r)
			if err != nil {
				return fmt.Errorf("invalid regex %q: %w", r, err)
			}

			var line int = -1
			offset := c.c.Value()
			for i := offset + 1; i != offset; i = (i + 1) % c.l.Len() {
				if regexp.MatchString(c.l.Index(i).String()) {
					line = i
					break
				}
			}

			if line == -1 {
				return c.errMsgf("No line matching regex %q found.\n", r)
			}

			err = c.c.Set(line)
			if err != nil {
				return err
			}
			return nil
		},
	}, {
		keys: []string{"goto", "g"},
		help: "Go to line number <N>.",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			n := args[0].(int)
			if l := c.l.Len(); n > l {
				return fmt.Errorf("line number too big: %d > %d", n, l)
			}

			if err := c.c.Set(n); err != nil {
				return err
			}

			return nil
		},
	}, {
		keys: []string{"entrypoint", "entry"},
		help: "Sets cursor to app entrypoint.",
		action: func(c *Control, args ...interface{}) error {
			a := c.p.Entrypoint()
			block, ok := c.p.Addr(a)
			if !ok {
				return fmt.Errorf("cannot find block at address 0x%x", a)
			}

			ins, ok := block.Addr(a)
			if !ok {
				return fmt.Errorf("cannot find instruction at address 0x%x", a)
			}

			line := c.l.Line(block, ins.Idx())
			c.c.Set(line)
			return nil
		},
	}, {
		keys: []string{"emulate", "emul", "e"},
		help: "Start emulating the machine code at current line.",
		action: func(c *Control, args ...interface{}) error {
			c.addMode("emulate", listEmulate())
			return nil
		},
	}, {
		keys: []string{"alllines"},
		help: "Prints all lines of the code into console. " +
			"Ignores current cursor position.",
		action: func(c *Control, args ...interface{}) error {
			for i := 0; i < c.l.Len(); i++ {
				fmt.Print(c.v.Format(i))
			}

			if err := c.errMsgf("\n"); err != nil {
				return err
			}

			return nil
		},
	}}
}

func listEmulate() []*command {
	return []*command{{}}
}
