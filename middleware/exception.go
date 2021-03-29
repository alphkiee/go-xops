package middleware

import (
	"go-xops/pkg/common"

	"github.com/gin-gonic/gin"
)

// 全局异常处理中间件
func Exception(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if resp, ok := err.(common.RespInfo); ok {
				common.JSON(c, common.Ok, resp)
				c.Abort()
				return
			}
			resp := common.RespInfo{
				Code:    common.InternalServerError,
				Data:    map[string]interface{}{},
				Message: common.CustomError[common.InternalServerError],
			}
			common.JSON(c, common.Ok, resp)
			c.Abort()
			return
		}
	}()
	c.Next()
}
