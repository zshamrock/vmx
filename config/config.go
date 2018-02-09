package config

import "os"

var configDir = os.ExpandEnv("$HOME/.vmx")

type Config struct {
	Dir string
}

var DefaultConfig = Config{configDir}
