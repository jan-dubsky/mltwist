package control

import (
	"decomp/internal/console/internal/cursor"
	"decomp/internal/console/internal/lines"
	"decomp/internal/console/internal/view"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Control struct {
	l *lines.Lines
	c *cursor.Cursor
	v *view.View

	reader   *lineReader
	commands map[string]*command
}

func New(l *lines.Lines, c *cursor.Cursor, v *view.View) *Control {
	return &Control{
		l: l,
		c: c,
		v: v,

		reader:   newLineReader(os.Stdin),
		commands: commandMap(),
	}
}

func (c *Control) parseCommand(str string) (*command, []interface{}, error) {
	parts := strings.Split(str, " ")
	cmdStr := parts[0]
	parts = parts[1:]

	cmd, ok := c.commands[cmdStr]
	if !ok {
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
	var cmdStr string
	for cmdStr == "" {
		fmt.Printf("Enter command: ")

		var err error
		cmdStr, err = c.reader.readLine()
		if err != nil {
			return err
		}
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
			return err
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
