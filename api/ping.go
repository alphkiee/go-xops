package api

import (
	"go-xops/pkg/common"

	"github.com/gin-gonic/gin"
)

// Ping doc
// @Summary Get Ping
// @Produce  json
// @Description 查看调用接口是否能够ping通
// @Success 200 {object} common.RespInfo
// @Router /api/ping [get]
func Ping(c *gin.Context) {
	common.Success()
}
