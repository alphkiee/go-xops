package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var ClientDocker *client.Client

// 初始化Docker Client
func InitClient() {
	clientDocker, err := client.NewClientWithOpts(client.WithHTTPClient(nil), client.WithHost("tcp://192.168.212.176:2375"), client.WithDialContext(nil), client.WithVersion("1.40"))
	if err != nil {
		log.Printf("clientDocker init is %v", err)
	}
	ClientDocker = clientDocker
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

// 镜像构建
func BuildImage(c *gin.Context) {
	go func() {
		log.Printf("构建开始 pid = %v\n", unix.Getpid())
		log.Printf("异步ID=%v\n", GoID())
		ctx, _ := context.WithTimeout(context.Background(), time.Second*120)
		//defer cancel()

		tar, err := archive.TarWithOptions("/Users/痞老板/Work/Golang/go-xops/build/rocketmq", &archive.TarOptions{})
		if err != nil {
			log.Printf("archive.TarWithOptions is %v", err)
		}
		opts := types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{"qperixdkajciospo" + "rocketmq_test"},
			Remove:     true,
		}
		res, err := ClientDocker.ImageBuild(ctx, tar, opts)
		if err != nil {
			log.Printf("Imagebuild is %v", err)
		}
		log.Printf("res.Body is = %v", res.Body)

		defer res.Body.Close()
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			lastLine := scanner.Text()
			fmt.Println(lastLine)
		}
	}()
}

func BuildImage2(c *gin.Context) {
	c.Copy()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*120)
	//defer cancel()

	tar, err := archive.TarWithOptions("/Users/痞老板/Work/Golang/go-xops/build/rocketmq", &archive.TarOptions{})
	if err != nil {
		log.Printf("archive.TarWithOptions is %v", err)
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{"qperixdkajciospo" + "rocketmq_test"},
		Remove:     true,
	}

	go func() {
		log.Printf("异步ID=%v\n", GoID())
		res, err := ClientDocker.ImageBuild(ctx, tar, opts)
		if err != nil {
			log.Printf("Imagebuild is %v", err)
		}
		log.Printf("res.Body is = %v", res.Body)

		defer res.Body.Close()
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			lastLine := scanner.Text()
			fmt.Println(lastLine)
		}
	}()
}

func main() {

	InitClient()
	r := gin.Default()
	r.GET("/", BuildImage)
	r.Run("127.0.0.1:8000")

	/*
		// 1.创建路由
		// 默认使用了2个中间件Logger(), Recovery()
		r := gin.Default()
		// 1.异步
		r.GET("/long_async", func(c *gin.Context) {
			// 需要搞一个副本
			copyContext := c.Copy()
			// 异步处理
			go func() {
				time.Sleep(3 * time.Second)
				log.Println("异步执行：" + copyContext.Request.URL.Path)
				log.Printf("异步ID = %v", GoID())
			}()
		})
		// 2.同步
		r.GET("/long_sync", func(c *gin.Context) {
			time.Sleep(3 * time.Second)
			log.Println("同步执行：" + c.Request.URL.Path)
		})
		r.Run(":8000")

	*/
}
