package command

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// CmdList lists available custom command
func CmdList(c *cli.Context) {
	names := GetCommandNames()
	for _, name := range names {
		fmt.Print(name)
		if commands[name].requiresConfirmation {
			fmt.Printf(" (%s)\n", commandNameConfirmationSuffix)
		} else {
			fmt.Println("")
		}
	}
}
