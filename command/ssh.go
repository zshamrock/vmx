package command

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"github.com/kevinburke/ssh_config"
	"path/filepath"
	"io/ioutil"
	"strings"
)

// SSH implements scp connection to the remote instance
func SSH(host, command string, ch chan int) {
	fmt.Fprintf(os.Stdout, "Running command %s on host %s\n", command, host)
	user := ssh_config.Get(host, "User")
	hostname := ssh_config.Get(host, "Hostname")
	identityFile := ssh_config.Get(host, "IdentityFile")
	var identityFilePath string
	if len(identityFile) == 0 || identityFile ==  "~/.ssh/identity" {
		identityFilePath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	} else {
		identityFilePath = os.ExpandEnv(strings.Replace(identityFile, "~", "${HOME}", -1))
	}
	pk, _ := ioutil.ReadFile(identityFilePath)
	signer, _ := ssh.ParsePrivateKey([]byte(pk))
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", hostname), config)
	if err != nil {
		log.Panicln("Failed to dial:", err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		log.Panicln("Failed to create session:", err.Error())
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	if err := session.Run(command); err != nil {
		log.Panicln("Failed to run:", err.Error())
	}
	fmt.Fprintf(os.Stdout, "Command completed on the host %s\n", host)
	ch <- 0
}
