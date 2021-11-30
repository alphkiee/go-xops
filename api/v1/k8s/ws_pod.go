package k8s

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	k8s "go-xops/internal/service/k8s"
	"go-xops/pkg/common"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"os"
)

func ExecPod(c *gin.Context) {
	// 初始化pod所在的corev1资源组，发送请求
	// PodExecOptions struct 包括Container stdout stdout  Command 等结构
	// scheme.ParameterCodec 应该是pod 的GVK （GroupVersion & Kind）之类的
	var req k8s.ExecPod
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
	}
	quest := common.ClientSet.CoreV1().RESTClient().Post().Resource(req.Resource).Name(req.Name).Namespace(req.Namespace).SubResource("exec").VersionedParams(&apiv1.PodExecOptions{
		Command: []string{"bash"},
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}, scheme.ParameterCodec)
	// remotecommand 主要实现了http 转 SPDY 添加X-Stream-Protocol-Version相关header 并发送请求
	exec, err := remotecommand.NewSPDYExecutor(common.Config, "POST", quest.URL())

	// 检查是不是终端
	if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
		fmt.Errorf("stdin/stdout should be terminal")
	}
	// 这个应该是处理Ctrl + C 这种特殊键位
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		common.FailWithMsg(err.Error())
	}
	defer terminal.Restore(0, oldState)

	// 用IO读写替换 os stdout
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  screen,
		Stdout: screen,
		Stderr: screen,
		Tty:    false,
	}); err != nil {
		common.FailWithMsg(err.Error())
		return
	}
}

func PodLog(c *gin.Context) {
	var req k8s.LogPod
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	quest := common.ClientSet.CoreV1().Pods(req.Namespace).GetLogs(req.Name, &apiv1.PodLogOptions{})
	res, err := quest.Stream(context.TODO())
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}

	defer res.Close()

	scanner := bufio.NewScanner(res)
	for scanner.Scan() {
		lastLine := scanner.Text()
		fmt.Println(lastLine)
	}
}
