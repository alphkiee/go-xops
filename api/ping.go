package api

import (
	"go-xops/dto/response"

	"github.com/gin-gonic/gin"
)

// 检查服务器是否通畅
func Ping(c *gin.Context) {
	response.Success()
}
