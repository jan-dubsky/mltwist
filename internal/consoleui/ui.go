package consoleui

import (
	"errors"
	"fmt"
	"mltwist/internal/consoleui/internal/linereader"
	"mltwist/internal/consoleui/internal/view"
	"strings"
)

type UI struct {
	modeStack []namedMode
}

func New(initialMode Mode) (*UI, error) {
	c := &UI{}
	if err := c.AddMode("app", initialMode); err != nil {
		return nil, fmt.Errorf("invalid initial mode: %w", err)
	}

	return c, nil
}

func (c *UI) mode() namedMode { return c.modeStack[len(c.modeStack)-1] }

func (c *UI) cmd(s string) (Command, bool) {
	cmd, ok := c.mode().cmdMap[s]
	return cmd, ok
}

func (c *UI) AddMode(name string, mode Mode) error {
	m, err := newMode(name, mode)
	if err != nil {
		return fmt.Errorf("cannot process mode %q: %w", name, err)
	}

	c.modeStack = append(c.modeStack, m)
	return nil
}

func (c *UI) quitMode() error {
	m := c.modeStack[len(c.modeStack)-1]
	c.modeStack = c.modeStack[:len(c.modeStack)-1]

	if len(c.modeStack) == 0 {
		fmt.Printf("leaving app\n")
		if _, err := linereader.ReadLine(); err != nil {
			return fmt.Errorf("readline error: %w", err)
		}
		return ErrQuit
	}

	fmt.Printf("leaving mode %s\n", m.name)
	if _, err := linereader.ReadLine(); err != nil {
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

func (c *UI) parseCommand(str string) (Command, []interface{}, error) {
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

func (c *UI) processCommand() error {
	cmdStr, err := linereader.ReadLine()
	if err != nil {
		return err
	}
	if cmdStr == "" {
		return nil
	}

	cmd, args, err := c.parseCommand(cmdStr)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		if _, err := linereader.ReadLine(); err != nil {
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
		if _, err := linereader.ReadLine(); err != nil {
			return fmt.Errorf("readline error: %w", err)
		}
		return nil
	}

	return nil
}

func (c *UI) Run() error {
	for {
		e := view.NewComposite(c.mode().mode.View(), commandPrompt{})
		err := view.Print(e)
		if err != nil {
			return fmt.Errorf("cannot print screen: %w", err)
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
