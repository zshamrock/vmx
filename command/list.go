package command

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// CmdList lists available custom command
func CmdList(c *cli.Context) {
	CheckUpdate(c)
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
