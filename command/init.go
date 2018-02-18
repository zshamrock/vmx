package command

import (
	"fmt"
	"github.com/zshamrock/vmx/config"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

const (
	CommandsConfigFileName   = "commands"
	HostsConfigFileName      = "hosts"
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

func init() {
	fmt.Println("Reading commands and hosts groups...")
	cfg := config.DefaultConfig
	commands = readCommands(cfg)
	hostsGroups = readHostsGroups(cfg)
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
