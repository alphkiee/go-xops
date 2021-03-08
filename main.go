package main

import (
	"context"
	"fmt"
	"go-xops/initialize"
	"go-xops/pkg/common"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Go-Xops
// @version 2.0
// @description Go-Xops swagger接口文档
// @contact.name pilaoban
// @contact.url https://github.com/jkuup
// @contact.email alphkiee@gmail.com
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:9000
func main() {
	// 初始化配置
	initialize.InitConfig()
	// 初始化路由
	r := initialize.Routers()
	//初始化数据库
	initialize.Mysql()
	//初始化Redis
	// initialize.Redis()
	// 初始校验器
	initialize.Validate()
	// 初始化Casbin
	initialize.Casbin()
	// 初始化创建上传目录
	_ = os.MkdirAll(common.Conf.Upload.SaveDir+"/avatar/", 644)
	//redis := cache.NewStringOperation()
	//fmt.Println(redis.Set("b","xx"))

	//是否初始化数据(慎用) $go-xops init
	/*
		if len(os.Args) > 1 {
			if os.Args[1] == "init" {
				initialize.InitData()
				fmt.Println("数据初始化成功!")
				os.Exit(1)
			}
		}
	*/

	// 关闭cache连接池
	// defer common.Redis.Close()

	// 启动服务器
	host := "0.0.0.0"
	port := common.Conf.System.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	go func() {
		// 加入pprof性能分析
		if err := http.ListenAndServe(":8005", nil); err != nil {
		}
	}()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
	}

}
