package command

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"github.com/zshamrock/vmx/config"
)

const commandNameConfirmationSuffix = "!"

type Command struct {
	name, command        string
	requiresConfirmation bool
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
		if name == "DEFAULT" {
			continue
		}
		requiresConfirmation := strings.HasSuffix(name, commandNameConfirmationSuffix)
		name = strings.TrimSuffix(name, commandNameConfirmationSuffix)
		commands[name] = Command{
			name,
			section.Key("command").String(),
			requiresConfirmation}
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
		if name == "DEFAULT" {
			continue
		}
		hosts[name] = section.KeyStrings()
	}
	return hosts
}

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	args := c.Args()
	host := args.Get(0)
	commandName := args.Get(1)
	command := getCommand(commandName, c)
	hosts := readHosts(config.DefaultConfig)
	target, ok := hosts[host]
	if !ok {
		// First then try whether host:children exists
		target, ok = hosts[host+":children"]
		if ok {
			children := make([]string, 0, len(target))
			for _, t := range target {
				children = append(children, hosts[t]...)
			}
			target = children
		} else {
			target = make([]string, 0, 1)
			target = append(target, host)
			fmt.Fprintf(os.Stdout, "%s: host group \"%s\" is not defined, interpret it as the ad-hoc host\n",
				c.App.Name, host)
		}
	}
	var confirmation string
	if command.requiresConfirmation {
		fmt.Fprintf(os.Stdout, "Confirm to run \"%s\" command on %v - yes/no or y/n: ", command.name, target)
		fmt.Scanln(&confirmation)
	}
	confirmation = strings.ToLower(confirmation)
	if confirmation != "yes" && confirmation != "y" {
		return
	}
	fmt.Fprintf(os.Stdout, "Running command: %s on %v\n", command.command, target)
	ch := make(chan int, len(target))
	for _, t := range target {
		go SSH(t, command.command, ch)
	}
	for i := 0; i < len(target); i++ {
		<-ch
	}
}

func getCommand(commandName string, c *cli.Context) Command {
	commands := readCommands(config.DefaultConfig)
	command, ok := commands[commandName]
	if !ok {
		adhocCommand := strings.Join(c.Args().Tail(), " ")
		fmt.Fprintf(os.Stdout, "%s: custom command \"%s\" is not defined, interpret it as the ad-hoc command: %s\n",
			c.App.Name, commandName, adhocCommand)
		command = Command{"ad-hoc", adhocCommand, false}
	}
	return command
}
