package initialize

import (
	"github.com/docker/docker/client"
	"go-xops/pkg/common"
)

func InitDockerApi() {
	c := common.Conf.DockerApi
	// docker reset api 地址
	b, err := client.NewClientWithOpts(client.WithHTTPClient(nil), client.WithHost(c.Url), client.WithDialContext(nil), client.WithVersion("1.40"))
	if err != nil {
		panic(err)
	}
	common.DockerClient = b
}
