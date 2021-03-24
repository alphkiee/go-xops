package cmdb

import (
	"go-xops/internal/response"
	"go-xops/internal/service"
	"go-xops/internal/service/cmd"
	"go-xops/pkg/utils"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const summary = `Summary:
All Executed Hosts: %d
Success:            %d
Failures:           %d
`

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func CmdExec(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		i   int
		res []interface{}
	)
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
	}

	ids := utils.Str2UintArr(c.Param("ids"))
	cmds := utils.Str2Arr(c.Param("cmds"))

	s := service.New()
	hosts, err := s.GetHostByIds(ids)
	hostNum := len(hosts)
	//fail需要优化，切片不够大，暂时只能给具体值
	fail := &cmd.Counter{Hosts: make([]string, hostNum)}
	for _, host := range hosts {
		port, _ := strconv.Atoi(host.Port)
		i++
		wg.Add(1)
		go func(ip string) {
			ssh := cmd.NewSSH(
				ip,
				host.User,
				host.Password,
				host.PrivateKey,
				"",
				port,
			)
			n := ssh.Ping()
			if n != 0 {
				fail.Incre(host.IP)
			}
			for _, v := range cmds {
				wg.Add(1)
				go func(v string) {
					_, rep := ssh.CmdExec(v)
					s := &response.CmdRep{
						IP:     ip,
						Cmd:    v,
						Stdout: rep,
					}
					res = append(res, s)
					wg.Done()
				}(v)
			}
			wg.Done()
		}(host.IP)

		if i == 30 {
			wg.Wait()
			i = 0
		}
	}

	wg.Wait()

	logrus.Printf(summary, hostNum, hostNum-fail.Data, fail.Data)
	response.SuccessWithData(res)
}

func Sftp(c *gin.Context) {
	var (
		wg sync.WaitGroup
		i  int
	)

	ids := utils.Str2UintArr(c.PostForm("ids"))
	srcPath := utils.Str2Arr(c.PostForm("srcPath"))
	dstPath := utils.Str2Arr(c.PostForm("dstPath"))

	s := service.New()
	hosts, err := s.GetHostByIds(ids)
	if err != nil {
		response.FailWithMsg("服务器内部错误")
	}
	for _, host := range hosts {
		port, _ := strconv.Atoi(host.Port)
		i++
		config := new(cmd.ClientConfig)
		config.CreateClient(host.IP, port, host.User, host.Password)
		for i, v := range srcPath {
			wg.Add(1)
			go func(ip string) {
				config.Upload(v, dstPath[i])
				wg.Done()
			}(host.IP)
		}
		if i == 30 {
			wg.Wait()
			i = 0
		}
	}
	wg.Wait()
	response.Success()

}

func SftpDow(c *gin.Context) {
	var (
		wg sync.WaitGroup
		i  int
	)

	ids := utils.Str2UintArr(c.PostForm("ids"))
	srcPath := c.PostForm("srcPath")
	dstPath := c.PostForm("dstPath")

	s := service.New()
	hosts, err := s.GetHostByIds(ids)
	if err != nil {
		response.FailWithMsg("服务器内部错误")
	}
	for _, host := range hosts {
		port, _ := strconv.Atoi(host.Port)
		i++
		config := new(cmd.ClientConfig)
		config.CreateClient(host.IP, port, host.User, host.Password)
		wg.Add(1)
		go func(ip string) {
			config.Download(srcPath, dstPath)
			wg.Done()
		}(host.IP)

		if i == 30 {
			wg.Wait()
			i = 0
		}
	}
	wg.Wait()
	response.Success()

}
