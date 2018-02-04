package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

// GlobalFlags used
var GlobalFlags = []cli.Flag{}

// Commands available
var Commands = []cli.Command{}

// CommandNotFound is used by cli to display an error message when unknown command is asked
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.\n", c.App.Name, command, c.App.Name, c.App.Name)
	cli.ShowAppHelp(c)
	os.Exit(2)
}
