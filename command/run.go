package command

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"github.com/zshamrock/vmx/config"
)

const (
	CommandsConfigFileName   = "commands"
	HostsConfigFileName      = "hosts"
	SectionCommandKeyName    = "command"
	SectionWorkingDirKeyName = "workingdir"

	defaultSectionName            = "DEFAULT"
	hostsGroupArgsIndex           = 0
	commandNameArgsIndex          = 1
	commandNameConfirmationSuffix = "!"
	hostsGroupChildrenSuffix      = ":children"
)

type Command struct {
	name, command, workingDir string
	requiresConfirmation      bool
}

var commands map[string]Command
var hostsGroups map[string][]string

func init() {
	fmt.Println("Reading commands and hosts groups...")
	cfg := config.DefaultConfig
	commands = readCommands(cfg)
	hostsGroups = readHostsGroups(cfg)
}

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	command := getCommand(c)
	hosts := getHosts(c)
	var confirmation string
	if command.requiresConfirmation {
		fmt.Fprintf(os.Stdout, "Confirm to run \"%s\" command on %v - yes/no or y/n: ", command.name, hosts)
		fmt.Scanln(&confirmation)
	}
	confirmation = strings.ToLower(confirmation)
	if command.requiresConfirmation && confirmation != "yes" && confirmation != "y" {
		return
	}
	cmd := command.command
	if command.workingDir != "" {
		cmd = fmt.Sprintf("cd %s && %s", command.workingDir, cmd)
	}
	fmt.Fprintf(os.Stdout, "Running command: %s from %s on %v\n", command.command, command.workingDir, hosts)
	ch := make(chan int, len(hosts))
	for _, host := range hosts {
		go SSH(host, cmd, ch)
	}
	for i := 0; i < len(hosts); i++ {
		<-ch
	}
}

func getCommand(c *cli.Context) Command {
	args := c.Args()
	commandName := args.Get(commandNameArgsIndex)
	command, ok := commands[commandName]
	if !ok {
		adhocCommand := strings.Join(c.Args().Tail(), " ")
		fmt.Fprintf(os.Stdout, "%s: custom command \"%s\" is not defined, interpret it as the ad-hoc command: %s\n",
			c.App.Name, commandName, adhocCommand)
		command = Command{"ad-hoc", adhocCommand, "", false}
	}
	return command
}

func readCommands(config config.Config) map[string]Command {
	commands := make(map[string]Command)
	cfg, err := ini.Load(config.Dir + "/" + CommandsConfigFileName)
	cfg.BlockMode = false
	if err != nil {
		os.Exit(1)
	}
	sections := cfg.Sections()
	// There is always DEFAULT section, so exclude that one from the commands capacity
	for _, section := range sections {
		name := section.Name()
		if name == defaultSectionName {
			continue
		}
		requiresConfirmation := strings.HasSuffix(name, commandNameConfirmationSuffix)
		name = strings.TrimSuffix(name, commandNameConfirmationSuffix)
		workingDir := ""
		if section.HasKey(SectionWorkingDirKeyName) {
			workingDir = section.Key(SectionWorkingDirKeyName).String()
		}
		commands[name] = Command{
			name,
			section.Key(SectionCommandKeyName).String(),
			workingDir,
			requiresConfirmation}
	}
	return commands
}

func getHosts(c *cli.Context) []string {
	args := c.Args()
	hostsGroup := args.Get(hostsGroupArgsIndex)
	hosts, ok := hostsGroups[hostsGroup]
	if !ok {
		// First then try whether host:children exists
		hosts, ok = hostsGroups[hostsGroup+hostsGroupChildrenSuffix]
		if ok {
			children := make([]string, 0, len(hosts))
			for _, group := range hosts {
				children = append(children, hostsGroups[group]...)
			}
			hosts = children
		} else {
			hosts = []string{hostsGroup}
			fmt.Fprintf(os.Stdout, "%s: hosts group \"%s\" is not defined, interpret it as the ad-hoc host\n",
				c.App.Name, hostsGroup)
		}
	}
	return hosts
}

func readHostsGroups(config config.Config) map[string][]string {
	groups := make(map[string][]string)
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, config.Dir+"/"+HostsConfigFileName)
	cfg.BlockMode = false
	if err != nil {
		os.Exit(1)
	}
	sections := cfg.Sections()
	// There is always DEFAULT section, so exclude that one from the commands capacity
	for _, section := range sections {
		name := section.Name()
		if name == defaultSectionName {
			continue
		}
		groups[name] = section.KeyStrings()
	}
	return groups
}
