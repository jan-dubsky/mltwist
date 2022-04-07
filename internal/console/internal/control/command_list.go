package control

import (
	"decomp/internal/console/internal/lines"
	"fmt"
	"math"
	"regexp"
)

// ErrQuit is returned by command in case UI termination is required.
var ErrQuit = fmt.Errorf("app exit required")

// commandsProxy is a helped variable containing the same array array (the very
// same instance) as commands variable.
//
// An unpleasant consequence of Go in-package initialization order is that there
// is no way how to write help command referring commands variable. The problem
// is that help command closure has to refer commands array and command array
// has to contains command for help message (including the help closure). Go is
// then not able to find any order in which those 2 values should be
// initialized.
//
// A natural solution would be to write a function (getter) which will return
// value of commands array. This way, there will be no closure and the plain
// (getter) function itself doesn't require any initialization, so there exists
// an initialization order without loops. Unfortunately, Go doesn't care about
// real dependencies. Instead it analyzes dependencies only based on lexical
// references. So if the help closure refers the getter function, the getter
// function refers the commands array and the help closure refers the getter, Go
// still understands this as a loop and doesn't compile this code.
//
// The only solution to this case seems to be to create another symbol which
// won't depend directly on commands array, but which will have the same
// content. To achieve this, we laverage go init function which is granted to
// run in the package initialization phase. Consequently we are certain that
// value of commandsProxy will be set before it will be used by the help
// closure. In other words, existence of this helper variable is just a fancy
// way how to bypass Go technical limitations, but it doesn't bring any
// performance, neither functional impact for the code.
var commandsProxy []*command

func init() { commandsProxy = commands() }

func commands() []*command {
	return []*command{{
		keys: []string{"help", "h"},
		help: "Print help of all commands",
		action: func(c *Control, _ ...interface{}) error {
			for _, cmd := range commandsProxy {
				fmt.Printf("%s\t(args: %d, additional_args: %t)\n",
					cmd.keysString(),
					len(cmd.args),
					cmd.optionalArgs != nil,
				)
				fmt.Print(format(cmd.help, 1, width))
				fmt.Printf("\n")
			}

			if err := c.errMsgf("\n"); err != nil {
				return err
			}

			return nil
		},
	}, {
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
	}, {
		keys: []string{"quit", "q"},
		help: "Quit the app.",
		action: func(c *Control, args ...interface{}) error {
			fmt.Printf("\n")
			return ErrQuit
		},
	}}
}

func commandMap() map[string]*command {
	m := make(map[string]*command)
	for i, command := range commands() {
		for _, k := range command.keys {
			if _, ok := m[k]; ok {
				panic(fmt.Sprintf(
					"duplicate command key %q at position %d", k, i,
				))
			}

			m[k] = command
		}
	}

	return m
}
