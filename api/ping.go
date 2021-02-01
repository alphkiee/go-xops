package api

import (
	"go-xops/dto/response"

	"github.com/gin-gonic/gin"
)

// Ping doc
// @Summary Get Ping
// @Produce  json
// @Description 查看调用接口是否能够ping通
// @Success 200 {object} response.RespInfo
// @Router /api/ping [get]
func Ping(c *gin.Context) {
	response.Success()
}
