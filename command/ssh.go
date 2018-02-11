package command

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"github.com/kevinburke/ssh_config"
	"path/filepath"
	"io/ioutil"
)

// SSH implements scp connection to the remote instance
func SSH(host, command string) {
	user := ssh_config.Get(host, "User")
	hostname := ssh_config.Get(host, "Hostname")
	pk, _ := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
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
}
