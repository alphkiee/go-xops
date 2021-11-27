package initialize

import (
	"github.com/mittwald/goharbor-client/v4/apiv1"
	"go-xops/pkg/common"
)

func InitHarborClient() {
	// 用户名和密码输入错误返回结果仍然为空
	harborClient, err := apiv1.NewRESTClientForHost(common.Conf.HarborApi.Url, common.Conf.HarborApi.User, common.Conf.HarborApi.Password)
	if err != nil {
		panic(err)
	}
	common.HarborClient = harborClient
}
