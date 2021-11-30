package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-xops/internal/service/docker"
	"go-xops/pkg/common"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: time.Duration(time.Second * 5),
	}
)

type PushImageReq struct {
	Name string `json:"name"`
}

// BuildImage doc
// @Summary Get /api/v1/docker/build
// @Description 构建镜像
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/docker/build [get]
func BuildImage(c *gin.Context) {
	//ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	//token := c.Request.Header.Get("x-token")
	//if token == "" {
	//	token = c.Query("x-token")
	//}
	//if err != nil {
	//common.FailWithMsg(err.Error())
	//	return
	//}
	//defer ws.Close()

	err := docker.BuildImage()
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}

	//err = ws.WriteMessage(websocket.TextMessage, b)
	//if err != nil {
	//	return
	//}
	common.Success()
}

// BuildImageSocket doc
// @Summary Get /api/v1/docker/build/socket
// @Description 构建镜像
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/docker/build/socket [get]
func BuildImageSocket(c *gin.Context) {
	var dockerRegistryUserID = ""
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	tar, err := archive.TarWithOptions("/Users/痞老板/Work/Golang/go-xops/build/rocketmq", &archive.TarOptions{})
	if err != nil {
		panic(err)
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{dockerRegistryUserID + "tcloud.hub/library/rocketmq"},
		Remove:     true,
	}
	// 调用镜像构建接口
	res, err := common.DockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	var lastLine string
	scanner := bufio.NewScanner(res.Body)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	token := c.Request.Header.Get("x-token")
	if token == "" {
		token = c.Query("x-token")
	}
	defer ws.Close()
	for scanner.Scan() {
		lastLine = scanner.Text()
		ws.WriteMessage(websocket.TextMessage, []byte(lastLine))
	}
}

// PushImage doc
// @Summary Post /api/v1/docker/push/image
// @Description 构建镜像
// @Produce json
// @Param name query string false "name"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/docker/push/image [post]
func PushImage(c *gin.Context) {
	var req PushImageReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	s, err := docker.ImagePush(req.Name)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.SuccessWithData(s)
}

// BuildImageSocketPush doc
// @Summary Get /api/v1/docker/build/socket
// @Description 构建镜像
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/docker/build/socket [get]
func BuildImageSocketPush(c *gin.Context) {
	go func() {
		var dockerRegistryUserID = ""
		ctx, _ := context.WithTimeout(context.Background(), time.Second*120)
		//defer cancel()
		tar, err := archive.TarWithOptions("/Users/痞老板/Work/Golang/go-xops/build/rocketmq", &archive.TarOptions{})
		if err != nil {
			panic(err)
		}

		opts := types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{dockerRegistryUserID + "tcloud.hub/library/rocketmq"},
			Remove:     true,
		}
		// 调用镜像构建接口
		res, err := common.DockerClient.ImageBuild(ctx, tar, opts)
		if err != nil {
			panic(err)
		}

		defer res.Body.Close()
		var lastLine string
		scanner := bufio.NewScanner(res.Body)
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		token := c.Request.Header.Get("x-token")
		if token == "" {
			token = c.Query("x-token")
		}
		defer ws.Close()
		for scanner.Scan() {
			lastLine = scanner.Text()
			ws.WriteMessage(websocket.TextMessage, []byte(lastLine))
		}
		// 推送新建的镜像到私有仓库
		authConfig := types.AuthConfig{Username: "admin", Password: "Harbor12345", ServerAddress: "tcloud.hub"}
		encodeJson, err := json.Marshal(authConfig)
		if err != nil {
			panic(err)
		}
		authStr := base64.URLEncoding.EncodeToString(encodeJson)
		// 调用官方docker reset api 接口推送镜像
		read, err := common.DockerClient.ImagePush(ctx, "tcloud.hub/library/rocketmq", types.ImagePushOptions{
			All:           false,
			RegistryAuth:  authStr,
			PrivilegeFunc: nil,
		})
		if err != nil {
			panic(err)
		}
		scanner = bufio.NewScanner(read)
		for scanner.Scan() {
			lastLine = scanner.Text()
			ws.WriteMessage(websocket.TextMessage, []byte(lastLine))
		}
	}()
}
