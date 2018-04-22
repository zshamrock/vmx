package config

import (
	"os"
	"path/filepath"
)

const (
	vmxHomeEnvVar          = "VMX_HOME"
	defaultVmxHome         = "${HOME}/.vmx"
	vmxSSHConfigHomeEnvVar = "VMX_SSH_CONFIG_HOME"
	defaultSSHConfigHome   = "${HOME}/.ssh"
)

type VMXConfig struct {
	Dir          string
	SSHConfigDir string
}

func (c VMXConfig) GetDir(profile string) string {
	if profile == "" {
		return c.Dir
	}
	return filepath.Join(c.Dir, profile)
}

var DefaultConfig VMXConfig

func init() {
	DefaultConfig = VMXConfig{os.ExpandEnv(defaultVmxHome), os.ExpandEnv(defaultSSHConfigHome)}
	vmxHome, ok := os.LookupEnv(vmxHomeEnvVar)
	if ok {
		DefaultConfig.Dir = vmxHome
	}
	vmxSSHHome, ok := os.LookupEnv(vmxSSHConfigHomeEnvVar)
	if ok {
		DefaultConfig.SSHConfigDir = vmxSSHHome
	}
}
