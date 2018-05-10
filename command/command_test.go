package command

import (
	"flag"
	"testing"

	"github.com/zshamrock/vmx/config"
	"github.com/zshamrock/vmx/core"
	"gopkg.in/urfave/cli.v1"
)

func TestGetCommand(t *testing.T) {
	followFlags := []string{"-f", "--follow"}
	for _, followFlag := range followFlags {
		flags := flag.FlagSet{}
		flags.Bool("follow", false, "")
		commandText := "tail -f -n 10 logs/rest.log"
		flags.Parse([]string{"--", followFlag, "dev", commandText})
		app := cli.NewApp()
		context := cli.NewContext(app, &flags, nil)
		command, extraArgs := getCommand(context, true)
		if !command.IsAdHoc() {
			t.Errorf("Command name should be ad-hoc, but got %s", command.Name)
		}
		if command.Command != commandText {
			t.Errorf("Command should be %s, but got %s", commandText, command.Command)
		}
		if extraArgs != "" {
			t.Errorf("Extra args should be empty, but got %s", extraArgs)
		}
	}
}

func TestGetCommandExtraArgs(t *testing.T) {
	config.Init("")
	followFlags := []string{"", "-f", "--follow"}
	for _, followFlag := range followFlags {
		commandText := "logs-extra"
		extraText := "rest.log"
		flags := flag.FlagSet{}
		follow := false
		arguments := []string{"dev", commandText, extraText}
		if followFlag != "" {
			follow = true
			flags.Bool("follow", false, "")
			arguments = append([]string{"--", followFlag}, arguments...)
		}
		flags.Parse(arguments)
		app := cli.NewApp()
		context := cli.NewContext(app, &flags, nil)
		command, extraArgs := getCommand(context, follow)
		expectedCommand := core.Command{
			Name:       "logs-extra",
			Command:    "tail -f -n 10 logs/%s",
			WorkingDir: "",
		}
		if command != expectedCommand {
			t.Errorf("Command should be %v, but got %v", expectedCommand, command)
		}
		if extraArgs != extraText {
			t.Errorf("Extra args should %s, but got %s", extraText, extraArgs)
		}
	}
}

func TestContainsFollow(t *testing.T) {
	followFlags := []string{"-f", "--follow"}
	for _, followFlag := range followFlags {
		flags := flag.FlagSet{}
		flags.Parse([]string{"--", followFlag, "dev", "tail -f -n 10 logs/rest.log"})
		app := cli.NewApp()
		context := cli.NewContext(app, &flags, nil)
		follow := ContainsFollow(context)
		if !follow {
			t.Error("Should contain follow")
		}
	}
}
