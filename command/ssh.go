package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

const (
	SshConfigUserKey         = "User"
	SshConfigHostnameKey     = "Hostname"
	SshConfigIdentityFileKey = "IdentityFile"
	ignoredIdentitySshFile   = "~/.ssh/identity"
)

// SSH implements scp connection to the remote instance
func SSH(sshConfig *ssh_config.Config, host, command string, ch chan int) {
	fmt.Printf("Running command: %s on host %s\n", command, host)
	user, _ := sshConfig.Get(host, SshConfigUserKey)
	hostname, _ := sshConfig.Get(host, SshConfigHostnameKey)
	identityFile, _ := sshConfig.Get(host, SshConfigIdentityFileKey)
	var identityFilePath string
	if identityFile == "" || identityFile == ignoredIdentitySshFile {
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
	var output strings.Builder
	session.Stdout = &output
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	if err := session.Run(command); err != nil {
		log.Panicln("Failed to run:", err.Error())
	}
	fmt.Fprintf(&output, "Command completed on the host %s\n", host)
	fmt.Println(output.String())
	ch <- 0
}
