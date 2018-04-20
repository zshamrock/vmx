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
	"gopkg.in/urfave/cli.v1"
)

const (
	hostsGroupArgsIndex      = 0
	commandNameArgsIndex     = 1
	hostsGroupChildrenSuffix = ":children"
	allHostsGroup            = "all"
)

// CmdRun runs custom command
func CmdRun(c *cli.Context) {
	CheckUpdate(c)
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
	sshConfig := readSSHConfig(config.DefaultConfig)
	ch := make(chan ExecOutput, len(hosts))
	for _, host := range hosts {
		if command.workingDir == "" && !strings.Contains(cmd, "cd ") {
			// Try to extend the command with the working dir from the defaults config, unless the command already has
			// have one, which takes the precedence. Also avoid to extend the command with the working dir from the
			// defaults config, if the command has "cd " in it, assuming user configured the working dir explicitly.
			defaults := getDefaults(host)
			workingDir, ok := defaults[SectionWorkingDirKeyName]
			if ok {
				fmt.Printf("Using working dir %s from the defaults config\n", workingDir)
				cmd = fmt.Sprintf("cd %s && %s", workingDir, cmd)
			}
		}
		go SSH(sshConfig, host, cmd, ch)
	}
	outputs := make([]ExecOutput, 0, len(hosts))
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
	hosts := getHostsByGroup(c, hostsGroup)
	sort.Strings(hosts)
	return hosts
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
				_, ok = hostsGroups[group+hostsGroupChildrenSuffix]
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

func getDefaults(host string) map[string]string {
	values := defaults[host]
	if values == nil {
		values = defaults[allHostsGroup]
	}
	if values == nil {
		return map[string]string{}
	}
	return values
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
