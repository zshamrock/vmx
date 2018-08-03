// This test requires VMX_HOME and VMX_SSH_CONFIG_HOME set to the test/config and test/ssh accordingly
package config

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/zshamrock/vmx/core"
)

func TestReadConfig(t *testing.T) {
	Init("")
	commands := GetCommands()
	expected := map[string]core.Command{
		"logs": {
			Name:       "logs",
			Command:    "cat logs/app.log",
			WorkingDir: "",
		},
		"app-logs": {
			Name:       "app-logs",
			Command:    "tail -f -n 10 logs/app.log",
			WorkingDir: "",
		},
		"follow-logs": {
			Name:       "follow-logs",
			Command:    "tail -f -n 10 logs/app1.log",
			WorkingDir: "/opt/app1",
			Follow:     true,
		},
		"redeploy": {
			Name:                 "redeploy",
			Command:              "./redeploy.sh",
			WorkingDir:           "",
			RequiresConfirmation: true,
		},
		"disk-space": {
			Name:       "disk-space",
			Command:    "df -h",
			WorkingDir: "",
		},
		"logs-extra1": {
			Name:       "logs-extra1",
			Command:    "tail -f -n 10 logs/%s",
			WorkingDir: "",
		},
		"logs-extra2": {
			Name:       "logs-extra2",
			Command:    "tail -f -n %s logs/app.log",
			WorkingDir: "",
		},
		"logs-extra3": {
			Name:       "logs-extra3",
			Command:    "tail -f logs/app.log",
			WorkingDir: "",
		},
	}
	if diff := deep.Equal(commands, expected); diff != nil {
		t.Error(diff)
	}
}
