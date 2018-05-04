package command

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/zshamrock/vmx/config"
	"github.com/zshamrock/vmx/core"
	"gopkg.in/urfave/cli.v1"
)

const (
	optionalFollowArgsIndex = 0
	hostsGroupArgsIndex     = 0
	commandNameArgsIndex    = 1
	FollowArgName           = "follow"
)

type execOutput struct {
	name, host string
	output     string
}

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	CheckUpdate(c)
	follow := ContainsFollow(c)
	command, extraArgs := getCommand(c, follow)
	if !follow && command.Follow {
		follow = command.Follow
	}
	hosts := getHosts(c, follow)
	var confirmation string
	if command.RequiresConfirmation {
		fmt.Printf("Confirm to run \"%s\" command on %v - yes/no or y/n: ", command.Name, hosts)
		fmt.Scanln(&confirmation)
	}
	confirmation = strings.ToLower(confirmation)
	if command.RequiresConfirmation && confirmation != "yes" && confirmation != "y" {
		return
	}
	cmd := command.Command
	if command.WorkingDir != "" {
		cmd = strings.TrimSpace(fmt.Sprintf("cd %s && %s %s", command.WorkingDir, cmd, extraArgs))
	}
	fmt.Printf("Running command: %s from %s on %v\n", command.Command, command.WorkingDir, hosts)
	sshConfig := readSSHConfig(config.DefaultConfig)
	ch := make(chan execOutput, len(hosts))
	for _, host := range hosts {
		if command.WorkingDir == "" && !strings.Contains(cmd, "cd ") {
			// Try to extend the command with the working dir from the defaults config, unless the command already has
			// have one, which takes the precedence. Also avoid to extend the command with the working dir from the
			// defaults config, if the command has "cd " in it, assuming user configured the working dir explicitly.
			defaults := config.GetDefaults(host)
			workingDir, ok := defaults[config.SectionWorkingDirKeyName]
			if ok {
				fmt.Printf("Using working dir %s from the defaults config\n", workingDir)
				cmd = fmt.Sprintf("cd %s && %s", workingDir, cmd)
			}
		}
		go ssh(sshConfig, host, cmd, follow, ch)
	}
	outputs := make([]execOutput, 0, len(hosts))
	for i := 0; i < len(hosts); i++ {
		outputs = append(outputs, <-ch)
	}
	sort.Slice(outputs, func(i, j int) bool {
		return outputs[i].host < outputs[j].host
	})
	for _, output := range outputs {
		fmt.Println(output.output)
	}
}
func getCommand(c *cli.Context, follow bool) (core.Command, string) {
	args := c.Args()
	actualCommandNameArgsIndex := getActualArgsIndex(commandNameArgsIndex, follow)
	commandName := strings.TrimSpace(args.Get(actualCommandNameArgsIndex))
	command, ok := config.GetCommands()[commandName]
	if !ok {
		adhocCommand := strings.Join(args[actualCommandNameArgsIndex:], " ")
		fmt.Printf("%s: custom command \"%s\" is not defined, interpret it as the ad-hoc command: %s\n",
			c.App.Name, commandName, adhocCommand)
		command = core.Command{
			Name:                 core.AdHocCommandName,
			Command:              adhocCommand,
			WorkingDir:           "",
			RequiresConfirmation: false,
		}
	}
	extraArgs := ""
	if ok && c.NArg() > 2 {
		extraArgsIndex := 1
		if follow {
			extraArgsIndex = 2
		}
		extraArgs = strings.Join(args.Tail()[extraArgsIndex:], " ")
	}
	return command, extraArgs
}
func getActualArgsIndex(argsIndex int, follow bool) int {
	actualArgsIndex := argsIndex
	if follow {
		actualArgsIndex = argsIndex + 1
	}
	return actualArgsIndex
}
func ContainsFollow(c *cli.Context) bool {
	follow := c.Args().Get(optionalFollowArgsIndex)
	return follow == "-f" || follow == fmt.Sprintf("--%s", FollowArgName)
}

func getHosts(c *cli.Context, follow bool) []string {
	args := c.Args()
	actualHostsGroupArgsIndex := getActualArgsIndex(hostsGroupArgsIndex, follow)
	hostsGroup := strings.TrimSpace(args.Get(actualHostsGroupArgsIndex))
	hosts := getHostsByGroup(c, hostsGroup)
	sort.Strings(hosts)
	return hosts
}

func getHostsByGroup(c *cli.Context, hostsGroup string) []string {
	hostsGroups := config.GetHostsGroups()
	if hostsGroup == config.AllHostsGroup {
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
		hosts, ok = hostsGroups[hostsGroup+config.HostsGroupChildrenSuffix]
		if ok {
			children := make([]string, 0, len(hosts))
			for _, group := range hosts {
				_, ok = hostsGroups[group+config.HostsGroupChildrenSuffix]
				if ok {
					children = append(children, getHostsByGroup(c, group)...)
				} else {
					children = append(children, hostsGroups[group]...)
				}
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

func readSSHConfig(cfg config.VMXConfig) *ssh_config.Config {
	sshConfigFilePath := filepath.Join(cfg.SSHConfigDir, "config")
	data, err := ioutil.ReadFile(sshConfigFilePath)
	if err != nil {
		fmt.Printf("Failed to read SSH config %s due to %v\n", sshConfigFilePath, err)
		os.Exit(1)
	}
	var buffer bytes.Buffer
	buffer.Write(data)
	sshConfig, err := ssh_config.Decode(&buffer)
	if err != nil {
		fmt.Printf("Failed to parse SSH config %s due to %v\n", sshConfigFilePath, err)
		os.Exit(1)
	}
	return sshConfig
}
