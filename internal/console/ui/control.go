package ui

import (
	"errors"
	"fmt"
	"mltwist/internal/console/internal/view"
	"os"
	"strings"
)

type Control struct {
	reader    *lineReader
	modeStack []namedMode
}

func New(initialMode Mode) (*Control, error) {
	c := &Control{
		reader: newLineReader(os.Stdin),
	}

	if err := c.AddMode("app", initialMode); err != nil {
		return nil, fmt.Errorf("invalid initial mode: %w", err)
	}

	return c, nil
}

func (c *Control) mode() namedMode { return c.modeStack[len(c.modeStack)-1] }

func (c *Control) cmd(s string) (Command, bool) {
	cmd, ok := c.mode().cmdMap[s]
	return cmd, ok
}

func (c *Control) AddMode(name string, mode Mode) error {
	m, err := newMode(name, mode)
	if err != nil {
		return fmt.Errorf("cannot process mode %q: %w", name, err)
	}

	c.modeStack = append(c.modeStack, m)
	return nil
}

func (c *Control) quitMode() error {
	m := c.modeStack[len(c.modeStack)-1]
	c.modeStack = c.modeStack[:len(c.modeStack)-1]

	if len(c.modeStack) == 0 {
		fmt.Printf("leaving app\n")
		if _, err := c.reader.readLine(); err != nil {
			return fmt.Errorf("readline error: %w", err)
		}
		return ErrQuit
	}

	fmt.Printf("leaving mode %s\n", m.name)
	if _, err := c.reader.readLine(); err != nil {
		return fmt.Errorf("readline error: %w", err)
	}

	return nil
}

func dropEmptyStrs(strs []string) []string {
	filtered := make([]string, 0, len(strs))
	for _, s := range strs {
		if len(s) != 0 {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func (c *Control) parseCommand(str string) (Command, []interface{}, error) {
	parts := dropEmptyStrs(strings.Split(str, " "))
	cmdStr := parts[0]
	parts = parts[1:]

	cmd, ok := c.cmd(cmdStr)
	if !ok {
		return Command{}, nil, fmt.Errorf("command %q not recognized", cmdStr)
	}

	if l := len(cmd.Args); len(parts) < l {
		err := fmt.Errorf("too few args: command %q requires %d args", cmdStr, l)
		return Command{}, nil, err
	}

	args := make([]interface{}, 0, len(parts))
	for i, parseF := range cmd.Args {
		val, err := parseF(parts[i])
		if err != nil {
			err = fmt.Errorf("cannot parse argument %d: %w", i, err)
			return Command{}, nil, err
		}

		args = append(args, val)
	}

	parts = parts[len(cmd.Args):]
	if len(parts) == 0 {
		return cmd, args, nil
	}

	vals, err := cmd.OptionalArgs(parts)
	if err != nil {
		return Command{}, nil, fmt.Errorf("cannot parse optional arguments: %w", err)
	}
	args = append(args, vals...)

	return cmd, args, nil
}

func (c *Control) processCommand() error {
	cmdStr, err := c.reader.readLine()
	if err != nil {
		return err
	}
	if cmdStr == "" {
		return nil
	}

	cmd, args, err := c.parseCommand(cmdStr)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		if _, err := c.reader.readLine(); err != nil {
			return fmt.Errorf("readline error: %w", err)
		}
		return nil
	}

	err = cmd.Action(c, args...)
	if err != nil {
		if errors.Is(err, ErrQuit) {
			return c.quitMode()
		}

		fmt.Printf("error: %s\n", err.Error())
		if _, err := c.reader.readLine(); err != nil {
			return fmt.Errorf("readline error: %w", err)
		}
		return nil
	}

	return nil
}

func (c *Control) Run() error {
	for {
		e := view.NewView(c.mode().mode.Element(), commandPrompt{})
		err := view.Print(e)
		if err != nil {
			return fmt.Errorf("cannot print UI elements: %w", err)
		}

		err = c.processCommand()
		if err != nil {
			if errors.Is(err, ErrQuit) {
				return nil
			}

			return fmt.Errorf("cannot process command: %w", err)
		}
	}
}

func (c *Control) ErrMsgf(pattern string, args ...interface{}) error {
	fmt.Printf(pattern, args...)
	fmt.Printf("Press ENTER to continue\n")
	if _, err := c.reader.readLine(); err != nil {
		return fmt.Errorf("readline error: %w", err)
	}
	return nil
}
