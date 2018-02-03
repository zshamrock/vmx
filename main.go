package main

import "fmt"

func main() {
	println(usage())
}

func usage() string {
	return fmt.Sprintf(`
vmx is a tool for interacting with cloud instances (like AWS EC2, for example) over SSH
`)
}
