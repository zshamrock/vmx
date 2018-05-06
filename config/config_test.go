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
	}
	if diff := deep.Equal(commands, expected); diff != nil {
		t.Error(diff)
	}
}
