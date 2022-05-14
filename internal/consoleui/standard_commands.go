package consoleui

import (
	"fmt"
	"mltwist/internal/consoleui/internal/linereader"
	"strings"
)

// ErrQuit is returned by command in case UI termination is required.
var ErrQuit = fmt.Errorf("app exit required")

var standardCmds = []Command{{
	Keys: []string{"quit", "q"},
	Help: "Quit the app.",
	Action: func(c *UI, args ...interface{}) error {
		fmt.Printf("\n")
		return ErrQuit
	}},
}

func addHelp(cmds []Command) []Command {
	helpCmd := Command{
		Keys: []string{"help", "h"},
		Help: "Print help of all commands",
	}
	cmds = append([]Command{helpCmd}, cmds...)

	cmds[0].Action = func(_ *UI, _ ...interface{}) error {
		for _, cmd := range cmds {
			fmt.Printf("%s\t(args: %d, additional_args: %t)\n",
				strings.Join(cmd.Keys, ", "),
				len(cmd.Args),
				cmd.OptionalArgs != nil,
			)
			fmt.Print(format(cmd.Help, 1, width))
			fmt.Printf("\n")
		}

		if err := linereader.ErrMsgf("\n"); err != nil {
			return err
		}

		return nil
	}

	return cmds
}

func addStandardCmds(cmds []Command) []Command {
	std := make([]Command, len(standardCmds), len(standardCmds)+len(cmds))
	copy(std, standardCmds)

	return addHelp(append(std, cmds...))
}
