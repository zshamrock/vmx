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
func SSH(sshConfig *ssh_config.Config, host, command string, ch chan ExecOutput) {
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
		log.Panicf("Failed to dial to the host %s: %v\n", host, err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		log.Panicf("Failed to create session for the host %s: %v\n", host, err.Error())
	}
	defer session.Close()
	var output strings.Builder
	session.Stdout = &output
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	if err := session.Run(command); err != nil {
		log.Panicf("Failed to run command \"%s\" on the host %s: %v\n", command, host, err.Error())
	}
	fmt.Fprintf(&output, "Command completed on the host %s\n", host)
	ch <- ExecOutput{
		command,
		host,
		output.String(),
	}
}
