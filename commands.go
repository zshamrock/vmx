package main

import (
	"fmt"
	"os"

	"github.com/zshamrock/vmx/command"
	"gopkg.in/urfave/cli.v1"
)

const profileArgName = "profile"

// GlobalFlags used
var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  fmt.Sprintf("%s, p", profileArgName),
		Usage: "profile to use to read hosts and commands from",
	},
}

func getProfile(c *cli.Context) string {
	profile := c.GlobalString(profileArgName)
	if profile == "" {
		profile = os.Getenv("VMX_DEFAULT_PROFILE")
	}
	return profile
}

// Commands available
var Commands = []cli.Command{
	{
		Name:  "run",
		Usage: "Run custom command",
		Description: `Example of usage is below:
    run logs    => run logs command defined in the ~/.vmx/commands`,
		Action: command.CmdRun,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  fmt.Sprintf("%s, f", command.FollowArgName),
				Usage: "flag indicates that the provided command will not exit, but will follow the output instead",
			},
		},
		SkipFlagParsing: true,
		BashComplete: func(c *cli.Context) {
			var names []string
			if c.NArg() == 0 || (c.NArg() == 1 && command.ContainsFollow(c)) {
				names = command.GetHostNames()
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
