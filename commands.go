package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"github.com/zshamrock/vmx/command"
	"github.com/kevinburke/ssh_config"
	"path/filepath"
	"sort"
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
		Action:          command.CmdRun,
		Flags:           []cli.Flag{},
		SkipFlagParsing: true,
		BashComplete: func(c *cli.Context) {
			names := make([]string, 0)
			if c.NArg() == 0 {
				f, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
				if err != nil {
					cfg, err := ssh_config.Decode(f)
					if err != nil {
						for _, host := range cfg.Hosts {
							names = append(names, host.String())
						}
					}
				}
				names = append(names, command.GetHostNames()...)
				sort.Strings(names)
			} else {
				names = command.GetCommandNames()
			}
			for _, name := range names {
				fmt.Println(name)
			}
		},
	},
	{
		Name:  "list",
		Usage: "List available custom commands",
		Description: `Example of usage is below:
    list    => list available custom commands defined in the ~/.vmx/commands`,
		Action: command.CmdList,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound is used by cli to display an error message when unknown command is asked
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.\n", c.App.Name, command, c.App.Name, c.App.Name)
	cli.ShowAppHelp(c)
	os.Exit(2)
}
