package command

import (
	"gopkg.in/urfave/cli.v1"
	"fmt"
	"sort"
)

// CmdList lists available custom command
func CmdList(c *cli.Context) {
	names := make([]string, 0, len(commands))
	for _, command := range commands {
		names = append(names, command.name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Print(name)
		if commands[name].requiresConfirmation {
			fmt.Printf(" (%s)\n", commandNameConfirmationSuffix)
		} else {
			fmt.Println("")
		}
	}
}
