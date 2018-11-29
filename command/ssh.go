package command

import (
	"fmt"
	"github.com/zshamrock/vmx/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevinburke/ssh_config"
	cryptoSSH "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	SshConfigUserKey         = "User"
	SshConfigHostnameKey     = "Hostname"
	SshConfigIdentityFileKey = "IdentityFile"
	ignoredIdentitySshFile   = "~/.ssh/identity"
	knownHostsFileName       = "known_hosts"
)

// ssh implements scp connection to the remote instance
func ssh(sshConfig *ssh_config.Config, host, command string, follow bool, ch chan execOutput) {
	fmt.Printf("Running command: %s on host %s\n", command, host)
	user, _ := sshConfig.Get(host, SshConfigUserKey)
	hostname, _ := sshConfig.Get(host, SshConfigHostnameKey)
	identityFile, _ := sshConfig.Get(host, SshConfigIdentityFileKey)
	var identityFilePath string
	if identityFile == "" || identityFile == ignoredIdentitySshFile {
		identityFilePath = filepath.Join(config.DefaultConfig.SSHConfigDir, "id_rsa")
	} else {
		identityFilePath = strings.Replace(identityFile, "~", config.DefaultConfig.SSHConfigDir, -1)
	}
	pk, _ := ioutil.ReadFile(identityFilePath)
	signer, _ := cryptoSSH.ParsePrivateKey([]byte(pk))
	hostKeyCallback := buildHostKeyCallback()
	config := &cryptoSSH.ClientConfig{
		User: user,
		Auth: []cryptoSSH.AuthMethod{
			cryptoSSH.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}
	client, err := cryptoSSH.Dial("tcp", fmt.Sprintf("%s:22", hostname), config)
	if err != nil {
		log.Panicf("Failed to dial to the host %s: %v\n", host, err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		log.Panicf("Failed to create session for the host %s: %v\n", host, err.Error())
	}
	defer session.Close()
	var output strings.Builder
	if follow {
		session.Stdout = os.Stdout
	} else {
		session.Stdout = &output
	}
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	if err := session.Run(command); err != nil {
		log.Panicf("Failed to run command \"%s\" on the host %s: %v\n", command, host, err.Error())
	}
	fmt.Fprintf(&output, "Command completed on the host %s\n", host)
	ch <- execOutput{
		command,
		host,
		output.String(),
	}
}

func buildHostKeyCallback() cryptoSSH.HostKeyCallback {
	configuredDefaultKnownHostsFile := filepath.Join(config.DefaultConfig.SSHConfigDir, knownHostsFileName)
	_, err := os.Stat(configuredDefaultKnownHostsFile)
	knownHostsFiles := make([]string, 0, 2)
	if err == nil {
		knownHostsFiles = append(knownHostsFiles, configuredDefaultKnownHostsFile)
	}
	if defaultKnownHostsFile := filepath.Join(os.ExpandEnv(config.DefaultSSHConfigHome), knownHostsFileName); configuredDefaultKnownHostsFile != defaultKnownHostsFile {
		knownHostsFiles = append(knownHostsFiles, defaultKnownHostsFile)
	}
	if len(knownHostsFiles) == 0 {
		fmt.Printf("No %s files are found\n", knownHostsFileName)
		os.Exit(1)
	}
	hostKeyCallback, err := knownhosts.New(knownHostsFiles...)
	if err != nil {
		log.Panicf("Failed to to read %s file %v\n", knownHostsFileName, err.Error())
	}
	return hostKeyCallback
}
