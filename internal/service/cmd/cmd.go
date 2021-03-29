package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Host                 string
	User                 string
	Password             string
	IdentityFile         string
	IdentityFilePassword string
	Port                 int
}

//连接的配置
type ClientConfig struct {
	Host       string       //ip
	Port       int          // 端口
	Username   string       //用户名
	Password   string       //密码
	sshClient  *ssh.Client  //ssh client
	sftpClient *sftp.Client //sftp client
}

type Counter struct {
	sync.Mutex
	Data  int
	Hosts []string
}

//CmdRep ...
type CmdRep struct {
	IP     string      `json:"ip"`
	Cmd    string      `json:"cmd"`
	Stdout interface{} `json:"stdout"`
}

func (c *Counter) Incre(host string) {
	c.Lock()
	defer c.Unlock()
	c.Hosts[c.Data] = host
	c.Data++
}

func (c *Counter) Decre() {
	c.Lock()
	defer c.Unlock()
	c.Data--
}

func NewSSH(host, user, password, identityFile, idfilepass string, port int) *SSH {
	return &SSH{
		Host:                 host,
		User:                 user,
		Password:             password,
		Port:                 port,
		IdentityFile:         identityFile,
		IdentityFilePassword: idfilepass,
	}
}

func (s *SSH) Connect() (*ssh.Session, error) {
	auths := []ssh.AuthMethod{
		ssh.Password(s.Password),
	}
	key, err := ioutil.ReadFile(s.IdentityFile)
	if err == nil {
		key, err := DecryptKey(key, []byte(s.IdentityFilePassword))
		if err != nil {
			goto M
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err == nil {
			auths = append(auths, ssh.PublicKeys(signer))
		}
	}
M:
	config := &ssh.ClientConfig{
		User:    s.User,
		Auth:    auths,
		Timeout: time.Second * 10,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	client, err := ssh.Dial("tcp", addr, config)

	if err != nil {
		return nil, err
	}

	return client.NewSession()
}

func (s *SSH) Run(cmd string) (stdout, stderr io.Reader, err error) {
	session, err := s.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()

	stdout, _ = session.StdoutPipe()
	stderr, _ = session.StderrPipe()

	err = session.Run(cmd)
	return stdout, stderr, err
}

func (s *SSH) Ping() int {
	_, err := s.Connect()
	if err != nil {
		log.Printf("Ping is Faild")
		return 1
	}
	return 0
}

func (s *SSH) CmdExec(cmd string) (int, string) {
	var (
		res []byte
	)
	// initial session
	session, err := s.Connect()
	if err != nil {
		//log.Printf("%s execute failed ~> \n\033[31m%s\033[0m\n\n", s.Host, err)
		return 1, ""
	}
	defer session.Close()

	// set stdout, stderr
	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()

	// run command
	err = session.Run(cmd)
	if err != nil {
		res, _ = ioutil.ReadAll(stderr)

		//log.Printf("%s execute failed ~> \n\033[31m%s\033[0m", s.Host, string(res))
		return 1, string(res)
	}
	res, _ = ioutil.ReadAll(stdout)
	// log.Printf("%s execute ok ~> \n\033[32m%s\033[0m", s.Host, string(res))

	return 0, string(res)
}

func DecryptKey(key, password []byte) ([]byte, error) {
	block, rest := pem.Decode(key)
	if len(rest) > 0 {
		return nil, errors.New(fmt.Sprintf("Decrypt key error: %s", string(rest)))
	}

	if x509.IsEncryptedPEMBlock(block) {
		der, err := x509.DecryptPEMBlock(block, password)
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: der}), nil
	}
	return key, nil
}

func (cliConf *ClientConfig) CreateClient(host string, port int, username, password string) {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password
	cliConf.Port = port

	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	if sshClient, err = ssh.Dial("tcp", addr, &config); err != nil {
		log.Fatalln("error occurred:", err)
	}
	cliConf.sshClient = sshClient

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Fatalln("error occurred:", err)
	}
	cliConf.sftpClient = sftpClient
}

func (cliConf *ClientConfig) RunShell(shell string) {
	var (
		session *ssh.Session
		err     error
	)

	//获取session，这个session是用来远程执行操作的
	if session, err = cliConf.sshClient.NewSession(); err != nil {
		log.Fatalln("error occurred:", err)
	}
	//执行shell
	if _, err := session.CombinedOutput(shell); err != nil {
		fmt.Println(shell)
		log.Fatalln("error occurred:", err)
	}
}

func (cliConf *ClientConfig) Upload(srcPath, dstPath string) {
	srcFile, _ := os.Open(srcPath)                   //本地
	dstFile, _ := cliConf.sftpClient.Create(dstPath) //远程
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("error occurred:", err)
			} else {
				break
			}
		}
		_, _ = dstFile.Write(buf[:n])
	}
}

func (cliConf *ClientConfig) Download(srcPath, dstPath string) {
	srcFile, _ := cliConf.sftpClient.Open(srcPath) //远程
	dstFile, _ := os.Create(dstPath)               //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Fatalln("error occurred", err)
	}
}
