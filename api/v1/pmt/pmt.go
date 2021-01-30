package pmt

import (
	"go-xops/dto/response"
	"go-xops/dto/service/pmt"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Pmt ...
func Pmt(c *gin.Context) {
	// 参数绑定
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	keys := utils.Str2Arr(c.Param("key"))
	jobs := utils.Str2Arr(c.Param("job"))

	// 开启多线程
	res, err := pmt.PrometheusAPIQuery_Test(keys, jobs)
	if err != nil {
		response.FailWithMsg("服务器内部错误")
	}
	response.SuccessWithData(res)
}
