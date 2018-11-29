package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/zshamrock/vmx/core"
	"gopkg.in/ini.v1"
)

const (
	CommandsConfigFileName   = "commands"
	HostsConfigFileName      = "hosts"
	DefaultsConfigFileName   = "defaults"
	SectionCommandKeyName    = "command"
	SectionWorkingDirKeyName = "workingdir"
	SectionFollowKeyName     = "follow"

	defaultSectionName = "DEFAULT"
)

var commands map[string]core.Command
var hostsGroups map[string][]string
var commandNames []string
var hostNames []string
var defaults map[string]map[string]string

func Init(profile string) {
	cfg := DefaultConfig
	commands = readCommands(cfg, profile)
	hostsGroups = readHostsGroups(cfg, profile)
	defaults = readDefaults(cfg, profile)
}

func readCommands(config VMXConfig, profile string) map[string]core.Command {
	commands := make(map[string]core.Command)
	cfg, err := ini.Load(filepath.Join(config.GetDir(profile), CommandsConfigFileName))
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
		requiresConfirmation := strings.HasSuffix(name, CommandNameConfirmationSuffix)
		name = strings.TrimSuffix(name, CommandNameConfirmationSuffix)
		workingDir := ""
		if section.HasKey(SectionWorkingDirKeyName) {
			workingDir = section.Key(SectionWorkingDirKeyName).String()
		}
		follow := false
		if section.HasKey(SectionFollowKeyName) {
			follow, _ = section.Key(SectionFollowKeyName).Bool()
		}
		commands[name] = core.Command{
			Name:                 name,
			Command:              section.Key(SectionCommandKeyName).String(),
			WorkingDir:           workingDir,
			Follow:               follow,
			RequiresConfirmation: requiresConfirmation,
		}
	}
	return commands
}

func readHostsGroups(config VMXConfig, profile string) map[string][]string {
	groups := make(map[string][]string)
	cfg, err := ini.LoadSources(
		ini.LoadOptions{AllowBooleanKeys: true},
		filepath.Join(config.GetDir(profile), HostsConfigFileName))
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
		sort.Strings(groups[name])
	}
	return groups
}

func readDefaults(config VMXConfig, profile string) map[string]map[string]string {
	defaults := make(map[string]map[string]string)
	cfg, err := ini.Load(filepath.Join(config.GetDir(profile), DefaultsConfigFileName))
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
		workingDir := ""
		if section.HasKey(SectionWorkingDirKeyName) {
			workingDir = section.Key(SectionWorkingDirKeyName).String()
		}
		values, ok := defaults[name]
		if !ok {
			values = make(map[string]string)
			defaults[name] = values
		}
		values[SectionWorkingDirKeyName] = workingDir
	}
	return defaults
}

func GetCommandNames() []string {
	if commandNames == nil {
		commandNames = make([]string, 0, len(commands))
		for _, command := range commands {
			commandNames = append(commandNames, command.Name)
		}
		sort.Strings(commandNames)
	}

	return commandNames
}

func GetCommands() map[string]core.Command {
	return commands
}

func GetHostNames() []string {
	if hostNames == nil {
		names := make(map[string]int)
		// Reading hosts from the ~/.vmx/hosts
		for group := range hostsGroups {
			var name string
			if strings.HasSuffix(group, HostsGroupChildrenSuffix) {
				name = strings.TrimSuffix(group, HostsGroupChildrenSuffix)
			} else {
				name = group
			}
			if names[name] == 0 {
				names[name] = 1
			}
		}

		// Reading hosts from ~/.ssh/config
		f, _ := os.Open(filepath.Join(DefaultConfig.SSHConfigDir, "config"))
		cfg, _ := ssh_config.Decode(f)
		for _, host := range cfg.Hosts {
			for _, pattern := range host.Patterns {
				v := pattern.String()
				if strings.ContainsAny(v, "*?") {
					continue
				}
				if names[v] == 0 {
					names[v] = 1
				}
			}
		}

		hostNames = make([]string, 0, len(names))
		for name := range names {
			hostNames = append(hostNames, name)
		}
		sort.Strings(hostNames)
	}
	return hostNames
}

func GetHostsGroups() map[string][]string {
	return hostsGroups
}

func GetDefaults(host string) map[string]string {
	values := defaults[host]
	if values == nil {
		values = defaults[AllHostsGroup]
	}
	if values == nil {
		return map[string]string{}
	}
	return values
}
