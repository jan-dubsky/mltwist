package control

import (
	"errors"
	"fmt"
	"mltwist/internal/console/internal/cursor"
	"mltwist/internal/console/internal/lines"
	"mltwist/internal/console/internal/view"
	"mltwist/internal/deps"
	"os"
	"strings"
)

type Control struct {
	p *deps.Program
	l *lines.Lines
	c *cursor.Cursor
	v *view.View

	reader    *lineReader
	modeStack []mode
}

func New(p *deps.Program, l *lines.Lines, c *cursor.Cursor, v *view.View) *Control {
	return &Control{
		p: p,
		l: l,
		c: c,
		v: v,

		reader:    newLineReader(os.Stdin),
		modeStack: []mode{newMode("app", listDisassemble())},
	}
}

func (c *Control) cmd(s string) *command {
	return c.modeStack[len(c.modeStack)-1].cmds[s]
}

func (c *Control) addMode(name string, cmds []*command) {
	c.modeStack = append(c.modeStack, newMode(name, cmds))
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

func (c *Control) parseCommand(str string) (*command, []interface{}, error) {
	parts := dropEmptyStrs(strings.Split(str, " "))
	cmdStr := parts[0]
	parts = parts[1:]

	cmd := c.cmd(cmdStr)
	if cmd == nil {
		return nil, nil, fmt.Errorf("command %q not recognized", cmdStr)
	}

	if l := len(cmd.args); len(parts) < l {
		err := fmt.Errorf("too few args: command %q requires %d args", cmdStr, l)
		return nil, nil, err
	}

	args := make([]interface{}, 0, len(parts))
	for i, parseF := range cmd.args {
		val, err := parseF(parts[i])
		if err != nil {
			return nil, nil, fmt.Errorf("cannot parse argument %d: %w", i, err)
		}

		args = append(args, val)
	}

	parts = parts[len(cmd.args):]
	if len(parts) == 0 {
		return cmd, args, nil
	}

	vals, err := cmd.optionalArgs(parts)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse optional arguments: %w", err)
	}
	args = append(args, vals...)

	return cmd, args, nil
}

func (c *Control) Command() error {
	fmt.Printf("Enter command: ")

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

	err = cmd.action(c, args...)
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

func (c *Control) errMsgf(pattern string, args ...interface{}) error {
	fmt.Printf(pattern, args...)
	fmt.Printf("Press ENTER to continue\n")
	if _, err := c.reader.readLine(); err != nil {
		return fmt.Errorf("readline error: %w", err)
	}
	return nil
}
