// This test requires VMX_HOME and VMX_SSH_CONFIG_HOME set to the test/config and test/ssh accordingly
package config

import (
	"reflect"
	"testing"

	"github.com/zshamrock/vmx/core"
)

func TestReadConfig(t *testing.T) {
	Init("")
	commands := GetCommands()
	expected := map[string]core.Command{
		"logs": {
			Name:       "logs",
			Command:    "cat logs/app.log",
			WorkingDir: "/opt/app",
		},
		"app-logs": {
			Name:       "app-logs",
			Command:    "tail -f -n 10 logs/app.log",
			WorkingDir: "/opt/app",
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
			WorkingDir:           "/opt/app",
			RequiresConfirmation: true,
		},
		"disk-space": {
			Name:       "disk-space",
			Command:    "command=df -h",
			WorkingDir: "/opt/app",
		},
	}
	equal := reflect.DeepEqual(commands, expected)
	if !equal {
		t.Errorf("Read commands %v from the 'commands' config file don't match expected commands %v",
			commands, expected)
	}
}
