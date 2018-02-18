package command

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

const (
	hostsGroupArgsIndex      = 0
	commandNameArgsIndex     = 1
	hostsGroupChildrenSuffix = ":children"
)

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	command := getCommand(c)
	hosts := getHosts(c)
	var confirmation string
	if command.requiresConfirmation {
		fmt.Printf("Confirm to run \"%s\" command on %v - yes/no or y/n: ", command.name, hosts)
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
	fmt.Printf("Running command: %s from %s on %v\n", command.command, command.workingDir, hosts)
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
		fmt.Printf("%s: custom command \"%s\" is not defined, interpret it as the ad-hoc command: %s\n",
			c.App.Name, commandName, adhocCommand)
		command = Command{"ad-hoc", adhocCommand, "", false}
	}
	return command
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
			fmt.Printf("%s: hosts group \"%s\" is not defined, interpret it as the ad-hoc host\n",
				c.App.Name, hostsGroup)
		}
	}
	return hosts
}
