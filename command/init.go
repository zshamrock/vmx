package command

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/zshamrock/vmx/config"
	"gopkg.in/ini.v1"
)

const (
	CommandsConfigFileName   = "commands"
	HostsConfigFileName      = "hosts"
	DefaultsConfigFileName   = "defaults"
	SectionCommandKeyName    = "command"
	SectionWorkingDirKeyName = "workingdir"

	defaultSectionName            = "DEFAULT"
	commandNameConfirmationSuffix = "!"
)

type Command struct {
	name, command, workingDir string
	requiresConfirmation      bool
}

var commands map[string]Command
var hostsGroups map[string][]string
var commandNames []string
var hostNames []string
var defaults map[string]map[string]string

func init() {
	cfg := config.DefaultConfig
	commands = readCommands(cfg)
	hostsGroups = readHostsGroups(cfg)
	defaults = readDefaults(cfg)
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

func readDefaults(config config.Config) map[string]map[string]string {
	defaults := make(map[string]map[string]string)
	cfg, err := ini.Load(config.Dir + "/" + DefaultsConfigFileName)
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
			commandNames = append(commandNames, command.name)
		}
		sort.Strings(commandNames)
	}

	return commandNames
}

func GetHostNames() []string {
	if hostNames == nil {
		names := make(map[string]int)
		// Reading hosts from the ~/.vmx/hosts
		for group := range hostsGroups {
			var name string
			if strings.HasSuffix(group, hostsGroupChildrenSuffix) {
				name = strings.TrimSuffix(group, hostsGroupChildrenSuffix)
			} else {
				name = group
			}
			if names[name] == 0 {
				names[name] = 1
			}
		}

		// Reading hosts from ~/.ssh/config
		f, _ := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
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
