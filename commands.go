package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"github.com/zshamrock/vmx/command"
)

// GlobalFlags used
var GlobalFlags = []cli.Flag{}

// Commands available
var Commands = []cli.Command{
	{
		Name:  "run",
		Usage: "Run custom command",
		Description: `Example of usage is below:
    run logs    => run logs command defined in the ~/.vmx/commands`,
		Action: command.CmdRun,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound is used by cli to display an error message when unknown command is asked
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.\n", c.App.Name, command, c.App.Name, c.App.Name)
	cli.ShowAppHelp(c)
	os.Exit(2)
}
