package control

import (
	"fmt"
	"strconv"
	"strings"
)

type argParseFunc func(s string) (interface{}, error)
type optArgParseFunc func(s []string) ([]interface{}, error)

type command struct {
	keys         []string
	help         string
	args         []argParseFunc
	optionalArgs optArgParseFunc
	action       func(c *Control, args ...interface{}) error
}

func (c command) keysString() string { return strings.Join(c.keys, ", ") }

// parseNum parses an integer parameter out of string and then validates that
// the value is in between min and max. Both min and max are inclusive
// boundaries.
func parseNum(min, max int) argParseFunc {
	return func(s string) (interface{}, error) {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value %q: %w", s, err)
		}

		if v < min {
			return nil, fmt.Errorf(
				"value is less than allowed minimum: %d < %d", v, min)
		}
		if v > max {
			return nil, fmt.Errorf(
				"value is greater than allowed maximum: %d > %d", v, max)
		}

		return v, nil
	}
}

func parseString(s string) (interface{}, error) {
	return s, nil
}

func joinOptStrings(strs []string) ([]interface{}, error) {
	return []interface{}{strings.Join(strs, "")}, nil
}

func helpCmd(cmds []*command) *command {
	helpCmd := &command{
		keys: []string{"help", "h"},
		help: "Print help of all commands",
	}
	cmds = append([]*command{helpCmd}, cmds...)

	helpCmd.action = func(c *Control, _ ...interface{}) error {
		for _, cmd := range cmds {
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
	}

	return helpCmd
}

func addHelp(cmds []*command) []*command {
	helpCmd := &command{
		keys: []string{"help", "h"},
		help: "Print help of all commands",
	}
	cmds = append([]*command{helpCmd}, cmds...)

	helpCmd.action = func(c *Control, _ ...interface{}) error {
		for _, cmd := range cmds {
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
	}

	return cmds
}

func addStandardCmds(cmds []*command) []*command {
	std := make([]*command, len(standardCmds), len(standardCmds)+len(cmds))
	copy(std, standardCmds)

	return addHelp(append(std, cmds...))
}

type cmdMap map[string]*command

func newCmdMap(cmds []*command) cmdMap {
	cmds = addStandardCmds(cmds)

	m := make(cmdMap)
	for i, cmd := range cmds {
		for _, k := range cmd.keys {
			if _, ok := m[k]; ok {
				panic(fmt.Sprintf(
					"duplicate command key %q at position %d", k, i,
				))
			}

			m[k] = cmd
		}
	}

	return m
}

type mode struct {
	name string
	cmds cmdMap
}

func newMode(name string, cmds []*command) mode {
	return mode{
		name: name,
		cmds: newCmdMap(cmds),
	}
}
