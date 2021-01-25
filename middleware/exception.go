package middleware

import (
	"go-xops/dto/response"

	"github.com/gin-gonic/gin"
)

// 全局异常处理中间件
func Exception(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 判断是否正常http响应结果放通
			if resp, ok := err.(response.RespInfo); ok {
				response.JSON(c, response.Ok, resp)
				c.Abort()
				return
			}
			// 服务器异常
			resp := response.RespInfo{
				Code:    response.InternalServerError,
				Data:    map[string]interface{}{},
				Message: response.CustomError[response.InternalServerError],
			}
			// 以json方式写入响应
			response.JSON(c, response.Ok, resp)
			c.Abort()
			return
		}
	}()
	c.Next()
}
