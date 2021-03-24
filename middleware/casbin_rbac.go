package middleware

import (
	v1 "go-xops/api/v1/system"
	"go-xops/assets/system"
	"go-xops/internal/response"
	"go-xops/pkg/common"
	"strings"

	"github.com/gin-gonic/gin"
)

// Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware(c *gin.Context) {
	// 获取当前登录用户
	user := v1.GetCurrentUserFromCache(c)
	// 当前登录用户的角色关键字作为casbin访问实体sub
	sub := user.(system.SysUser).Role.Keyword
	// 请求URL路径作为casbin访问资源obj(需先清除path前缀)
	obj := strings.Replace(c.Request.URL.Path, "/"+common.Conf.System.UrlPathPrefix, "", 1)
	// 请求方式作为casbin访问动作act
	act := c.Request.Method
	// 获取casbin策略管理器
	e := common.Casbin
	// 检查策略
	pass, _ := e.Enforce(sub, obj, act)
	if !pass {
		response.FailWithCode(response.Forbidden)
	}
	// 处理请求
	c.Next()
}
