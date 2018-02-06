package config

import (
	"os"
	"gopkg.in/ini.v1"
)

var configDir = os.ExpandEnv("$HOME/.vmx")

type Command struct {
	name, command string
}

func ReadCommands() []Command {
	cfg, err := ini.Load(configDir + "/commands")
	if err != nil {
		os.Exit(1)
	}
	sections := cfg.Sections()
	commands := make([]Command, len(sections))
	for _, section := range sections {
		commands = append(commands, Command{section.Name(), section.Body()})
	}
	return commands
}
