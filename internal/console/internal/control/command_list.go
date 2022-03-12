package control

import (
	"decomp/internal/console/internal/lines"
	"fmt"
	"math"
)

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

func init() { commandsProxy = commands }

var commands = []*command{
	{
		keys: []string{"help", "h"},
		help: "Print help of all commands",
		action: func(c *Control, _ ...interface{}) error {
			for _, cmd := range commandsProxy {
				fmt.Printf("%s\n", cmd.keysString())
				fmt.Printf("%s\n", format(cmd.help, 1, width))
				fmt.Printf("\n")
			}

			fmt.Printf("Press ENTER to continue\n")
			if _, err := c.reader.readLine(); err != nil {
				return fmt.Errorf("readline error: %w", err)
			}

			return nil
		},
	},
	{
		keys: []string{"down", "d"},
		help: "Move line cursor <N> lines down.",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			return c.l.Shift(args[0].(int))
		},
	},
	{
		keys: []string{"up", "u"},
		help: "Move line cursor <N> lines up.",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			return c.l.Shift(-args[0].(int))
		},
	},
	{
		keys: []string{"move", "mv", "m"},
		help: "Move instruction from line <N> to line <M>",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			from, to := args[0].(int), args[1].(int)

			fromBlock, fromIns, err := insLine(c.l, from)
			if err != nil {
				return fmt.Errorf("invalid from: %w", err)
			}
			toBlock, toIns, err := insLine(c.l, to)
			if err != nil {
				return fmt.Errorf("invalid to: %w", err)
			}

			if fromBlock != toBlock {
				return fmt.Errorf("instructions cannot be moved in between blocks")
			}
			if err := fromBlock.Move(fromIns.Idx(), toIns.Idx()); err != nil {
				return err
			}

			c.l.Reload(fromBlock.Idx())

			c.l.UnmarkAll()
			c.l.SetMark(from, lines.MarkMovedFrom)
			c.l.SetMark(to, lines.MarkMovedTo)

			return nil
		},
	},
	{
		keys: []string{"bounds", "b"},
		help: "Show bounds where a given instruction can be moved.",
		args: []argParseFunc{
			parseNum(0, math.MaxInt),
		},
		action: func(c *Control, args ...interface{}) error {
			block, ins, err := insLine(c.l, args[0].(int))
			if err != nil {
				return err
			}

			lower := block.LowerBound(ins)
			upper := block.UpperBound(ins)
			lowerLine := c.l.Line(block, block.Index(lower))
			upperLine := c.l.Line(block, block.Index(upper))

			c.l.UnmarkAll()
			// Lower and Upper indices are inclusive, but in
			// visualization we want to have exclusive indices.
			c.l.SetMark(lowerLine-1, lines.MarkLowerBound)
			c.l.SetMark(upperLine+1, lines.MarkUpperBound)

			return nil
		},
	},
	{
		keys: []string{"quit", "q"},
		help: "Quit the app.",
		action: func(c *Control, args ...interface{}) error {
			return nil
		},
	},
}

var cmdMap = func() map[string]*command {
	m := make(map[string]*command)
	for i, command := range commands {
		for _, k := range command.keys {
			if _, ok := m[k]; ok {
				panic(fmt.Sprintf("duplicate command key %q at position %d",
					k, i))
			}

			m[k] = command
		}
	}

	return m
}()
