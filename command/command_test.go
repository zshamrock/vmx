package command

import (
	"flag"
	"testing"

	"github.com/zshamrock/vmx/config"
	"github.com/zshamrock/vmx/core"
	"gopkg.in/urfave/cli.v1"
)

func TestGetCommand(t *testing.T) {
	followFlags := []string{"", "-f", "--follow"}
	for _, followFlag := range followFlags {
		commandText := "tail -f -n 10 logs/rest.log"
		arguments := []string{"dev", commandText}
		flags := flag.FlagSet{}
		follow := false
		if followFlag != "" {
			flags.Bool("follow", false, "")
			follow = true
			arguments = append([]string{"--", followFlag}, arguments...)
		}
		flags.Parse(arguments)
		app := cli.NewApp()
		context := cli.NewContext(app, &flags, nil)
		command := getCommand(context, follow)
		if !command.IsAdHoc() {
			t.Errorf("Command name should be ad-hoc, but got %s", command.Name)
		}
		if command.Command != commandText {
			t.Errorf("Command should be %s, but got %s", commandText, command.Command)
		}
	}
}

func TestGetCommandExtraArgs(t *testing.T) {
	config.Init("")
	followFlags := []string{"", "-f", "--follow"}
	for _, followFlag := range followFlags {
		commandText := "logs-extra"
		extraText := "rest.log"
		arguments := []string{"dev", commandText, extraText}
		flags := flag.FlagSet{}
		follow := false
		if followFlag != "" {
			flags.Bool("follow", false, "")
			follow = true
			arguments = append([]string{"--", followFlag}, arguments...)
		}
		flags.Parse(arguments)
		app := cli.NewApp()
		context := cli.NewContext(app, &flags, nil)
		command := getCommand(context, follow)
		expectedCommand := core.Command{
			Name:       "logs-extra",
			Command:    "tail -f -n 10 logs/" + extraText,
			WorkingDir: "",
		}
		if command != expectedCommand {
			t.Errorf("Command should be %v, but got %v", expectedCommand, command)
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
