package command

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"github.com/zshamrock/vmx/config"
)

type Command struct {
	name, command string
}

func readCommands(config config.Config) map[string]Command {
	commands := make(map[string]Command)
	cfg, err := ini.Load(config.Dir + "/commands")
	cfg.BlockMode = false
	if err != nil {
		os.Exit(1)
	}
	sections := cfg.Sections()
	// There is always DEFAULT section, so exclude that one from the commands capacity
	for _, section := range sections {
		name := section.Name()
		if strings.Compare(name, "DEFAULT") == 0 {
			continue
		}
		commands[section.Name()] = Command{section.Name(), section.Key("command").String()}
	}
	return commands
}

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	commands := readCommands(config.DefaultConfig)
	commandName := c.Args().First()
	command, ok := commands[commandName]
	if !ok {
		names := make([]string, len(commands))
		i := 0
		for name := range commands {
			names[i] = name
			i++
		}
		fmt.Fprintf(os.Stderr, "%s: custom command \"%s\" is not defined\n", c.App.Name, commandName)
		fmt.Fprintf(os.Stdout, "Known list of custom commands are: %s\n", strings.Join(names, ", "))
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Running command: %s\n", command.command)
}
