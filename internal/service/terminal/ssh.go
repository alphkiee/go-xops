package terminal

import (
	"go-xops/assets/cmdb"
	"io/ioutil"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// NewSSHClient ...
func NewSSHClient(host *cmdb.Host) (*ssh.Client, error) {

	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		config       ssh.Config
		client       *ssh.Client
	)

	auth = make([]ssh.AuthMethod, 0)

	if host.AuthType == "密码认证" {
		auth = append(auth, ssh.Password(host.Password))
	} else {
		pemBytes, err := ioutil.ReadFile(host.PrivateKey)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer
		if host.Password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(host.Password))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	clientConfig = &ssh.ClientConfig{
		User:    host.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", host.IP, clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil

}
