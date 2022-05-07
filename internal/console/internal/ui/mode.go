package ui

import (
	"fmt"
	"mltwist/internal/console/internal/view"
)

type commandMap map[string]Command

func newCmdMap(cmds []Command) (commandMap, error) {
	cmds = addStandardCmds(cmds)

	m := make(commandMap, len(cmds))
	for i, cmd := range cmds {
		for _, k := range cmd.Keys {
			if _, ok := m[k]; ok {
				return nil, fmt.Errorf(
					"duplicate command key %q at position %d", k, i)
			}

			m[k] = cmd
		}
	}

	return m, nil
}

type Mode interface {
	Commands() []Command
	View() view.View
}

type namedMode struct {
	name   string
	mode   Mode
	cmdMap commandMap
}

func newMode(name string, mode Mode) (namedMode, error) {
	cmdMap, err := newCmdMap(mode.Commands())
	if err != nil {
		return namedMode{}, err
	}

	return namedMode{
		name:   name,
		mode:   mode,
		cmdMap: cmdMap,
	}, nil
}
