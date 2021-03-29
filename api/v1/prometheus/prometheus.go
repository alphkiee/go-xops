package prometheus

import (
	"go-xops/internal/service/prometheus"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Pmt doc
// @Summarg Get /api/v1/prometheus/host/:key/:job
// @Description 根据key, job 获取监控项值
// @Produce json
// @Param key path string true "key"
// @Param job path string true "job"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 500 {object} common.RespInfo
// @Router /api/v1/prometheus/host/:key/:job [get]
func Pmt(c *gin.Context) {
	// 参数绑定
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	keys := utils.Str2Arr(c.Param("key"))
	jobs := utils.Str2Arr(c.Param("job"))

	// 开启多线程
	res, err := prometheus.PrometheusAPIQuery_Test(keys, jobs)
	if err != nil {
		common.FailWithMsg("服务器内部错误")
	}
	common.SuccessWithData(res)
}
