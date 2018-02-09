package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Authors = []cli.Author{{Name: "Aliaksandr Kazlou"}}
	app.Usage = usage()

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}

func usage() string {
	return fmt.Sprintf(`
vmx is a tool for interacting with cloud instances (like AWS EC2, for example) over SSH
`)
}
