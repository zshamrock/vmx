package main

import (
	"fmt"
	"os"

	"github.com/zshamrock/vmx/config"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Authors = []cli.Author{{Name: "Aliaksandr Kazlou", Email: "aliaksandr.kazlou@gmail.com"}}
	app.Metadata = map[string]interface{}{"GitHub": "https://github.com/zshamrock/vmx"}
	app.Usage = usage()
	app.EnableBashCompletion = true

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	app.Before = func(c *cli.Context) error {
		profile := getProfile(c)
		config.Init(profile)
		return nil
	}
	app.Run(os.Args)
}

func usage() string {
	return fmt.Sprintf(`
vmx is a tool for interacting with cloud instances (like AWS EC2, for example) over SSH
[https://github.com/zshamrock/vmx]
`)
}
