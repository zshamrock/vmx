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
	allHostsGroup            = "all"
)

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	command, extraArgs := getCommand(c)
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
		cmd = strings.TrimSpace(fmt.Sprintf("cd %s && %s %s", command.workingDir, cmd, extraArgs))
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

func getCommand(c *cli.Context) (Command, string) {
	args := c.Args()
	commandName := strings.TrimSpace(args.Get(commandNameArgsIndex))
	command, ok := commands[commandName]
	if !ok {
		adhocCommand := strings.Join(c.Args().Tail(), " ")
		fmt.Printf("%s: custom command \"%s\" is not defined, interpret it as the ad-hoc command: %s\n",
			c.App.Name, commandName, adhocCommand)
		command = Command{"ad-hoc", adhocCommand, "", false}
	}
	extraArgs := ""
	if ok && c.NArg() > 2 {
		extraArgs = strings.Join(c.Args().Tail()[1:], " ")
	}
	return command, extraArgs
}

func getHosts(c *cli.Context) []string {
	args := c.Args()
	hostsGroup := strings.TrimSpace(args.Get(hostsGroupArgsIndex))
	return getHostsByGroup(c, hostsGroup)
}

func getHostsByGroup(c *cli.Context, hostsGroup string) []string {
	if hostsGroup == allHostsGroup {
		allHosts := make([]string, 0, len(hostsGroups))
		for _, hosts := range hostsGroups {
			for _, host := range hosts {
				allHosts = append(allHosts, getHostsByGroup(c, host)...)
			}
		}
		return allHosts
	}
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
