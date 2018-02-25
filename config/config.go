package config

import "os"

const (
	vmxHomeEnvVar  = "VMXHOME"
	defaultVmxHome = "${HOME}/.vmx"
)

type Config struct {
	Dir string
}

var DefaultConfig Config

func init() {
	DefaultConfig = Config{os.ExpandEnv(defaultVmxHome)}
	vmxhome, ok := os.LookupEnv(vmxHomeEnvVar)
	if ok {
		DefaultConfig = Config{vmxhome}
	}
}
