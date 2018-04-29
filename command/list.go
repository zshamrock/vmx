package command

import (
	"fmt"

	"github.com/zshamrock/vmx/config"
	"gopkg.in/urfave/cli.v1"
)

// CmdList lists available custom command
func CmdList(c *cli.Context) {
	CheckUpdate(c)
	names := config.GetCommandNames()
	commands := config.GetCommands()
	for _, name := range names {
		fmt.Print(name)
		if commands[name].RequiresConfirmation {
			fmt.Printf(" (%s)\n", config.CommandNameConfirmationSuffix)
		} else {
			fmt.Println("")
		}
	}
}
