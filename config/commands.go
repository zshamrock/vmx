package config

import (
	"os"
	"gopkg.in/ini.v1"
	"strings"
)

var configDir = os.ExpandEnv("$HOME/.vmx")

type Command struct {
	name, command string
}

func ReadCommands() []Command {
	cfg, err := ini.Load(configDir + "/commands")
	cfg.BlockMode = false
	if err != nil {
		os.Exit(1)
	}
	sections := cfg.Sections()
	// There is always DEFAULT section, so exclude that one from the commands capacity
	commands := make([]Command, 0, len(sections)-1)
	for _, section := range sections {
		name := section.Name()
		if strings.Compare(name, "DEFAULT") == 0 {
			continue
		}
		commands = append(commands, Command{section.Name(), section.Key("command").String()})
	}
	return commands
}
