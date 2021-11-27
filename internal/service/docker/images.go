package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"go-xops/assets/docker"
	"go-xops/pkg/common"
	"io"
	"time"
)

type Images struct {
	I []docker.Image
}

var dockerRegistryUserID = ""

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func GetImageList(ctx context.Context) (Images, error) {
	var Is Images
	var I docker.Image
	images, err := common.DockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return Is, err
	}

	for _, image := range images {
		I.ImageID = image.ID
		I.Created = image.Created
		I.Size = image.Size
		I.Tag = image.RepoTags
		I.Repository = image.RepoDigests
		Is.I = append(Is.I, I)
	}
	return Is, nil
}

func imageBuild() error {
	//var bt []byte
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions("/Users/痞老板/Work/Golang/go-xops/build/py", &archive.TarOptions{})
	if err != nil {
		return err
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{dockerRegistryUserID + "py_xops"},
		Remove:     true,
	}
	res, err := common.DockerClient.ImageBuild(ctx, tar, opts)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	buildLog(res.Body)
	//bt, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func buildLog(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// 构建镜像(目前是把目录文件固定)
func BuildImage() error {
	err := imageBuild()
	if err != nil {
		return err
	}
	return nil
}

// 推送镜像到私有仓库
func ImagePush(name string) (string, error) {
	var rep string
	authConfig := types.AuthConfig{Username: "admin", Password: "Harbor12345", ServerAddress: "tcloud.hub"}
	encodeJson, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
		return rep, err
	}
	authStr := base64.URLEncoding.EncodeToString(encodeJson)
	// 调用官方docker reset api 接口推送镜像
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	readCloser, err := common.DockerClient.ImagePush(ctx, name, types.ImagePushOptions{
		All:           false,
		RegistryAuth:  authStr,
		PrivilegeFunc: nil,
	})
	if err != nil {
		panic(err)
		return rep, err
	}
	defer readCloser.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(readCloser)
	rep = buf.String()
	return rep, nil
}
