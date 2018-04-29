package command

import (
	"flag"
	"testing"

	"gopkg.in/urfave/cli.v1"
)

func TestGetCommand(t *testing.T) {
	flags := flag.FlagSet{}
	flags.Bool("follow", false, "")
	commandText := "tail -f -n 10 logs/rest.log"
	flags.Parse([]string{"--", "--follow", "dev", commandText})
	app := cli.NewApp()
	context := cli.NewContext(app, &flags, nil)
	command, extraArgs := getCommand(context, true)
	if !command.IsAdHoc() {
		t.Errorf("Command name should be ad-hoc, but got %s", command.name)
	}
	if command.command != commandText {
		t.Errorf("Command should be %s, but got %s", commandText, command.command)
	}
	if extraArgs != "" {
		t.Errorf("Extra args should be empty, but got %s", extraArgs)
	}
}

func TestContainsFollow(t *testing.T) {
	flags := flag.FlagSet{}
	flags.Bool("follow", false, "")
	flags.Parse([]string{"--", "--follow", "dev", "tail -f -n 10 logs/rest.log"})
	app := cli.NewApp()
	context := cli.NewContext(app, &flags, nil)
	follow := ContainsFollow(context)
	if !follow {
		t.Error("Should contain follow")
	}
}
