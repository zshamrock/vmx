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
		commands[name] = Command{name, section.Key("command").String()}
	}
	return commands
}

func readHosts(config config.Config) map[string][]string {
	hosts := make(map[string][]string)
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, config.Dir+"/hosts")
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
		hosts[name] = section.KeyStrings()
	}
	return hosts
}

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	commands := readCommands(config.DefaultConfig)
	args := c.Args()
	host := args.Get(0)
	commandName := args.Get(1)
	command, ok := commands[commandName]
	if !ok {
		names := make([]string, len(commands))
		i := 0
		for name := range commands {
			names[i] = name
			i++
		}
		adhocCommand := strings.Join(args.Tail(), " ")
		fmt.Fprintf(os.Stdout, "%s: custom command \"%s\" is not defined, interpret it as the ad-hoc command: %s\n",
			c.App.Name, commandName, adhocCommand)
		command = Command{"ad-hoc", adhocCommand}
	}
	hosts := readHosts(config.DefaultConfig)
	target := hosts[host]
	fmt.Fprintf(os.Stdout, "Running command: %s on %v\n", command.command, target)
	SSH(target[0], command.command)
}
