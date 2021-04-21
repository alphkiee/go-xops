package initialize

import (
	"fmt"
	"go-xops/internal/routers"
	"go-xops/internal/routers/cmdb"
	"go-xops/internal/routers/k8s"
	"go-xops/internal/routers/prometheus"
	"go-xops/internal/routers/system"
	"go-xops/middleware"
	"go-xops/pkg/common"

	_ "go-xops/docs"

	"github.com/gin-gonic/gin"
	swagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Routers ...
func Routers() *gin.Engine {
	gin.SetMode(common.Conf.System.AppMode)
	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	// r := gin.Default()
	// 创建不带中间件的路由:
	r := gin.New()
	// 添加中间件logger记录日志
	r.Use(middleware.LoggerToFile())
	r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	// 添加全局异常处理中间件
	r.Use(middleware.Exception)
	// 添加跨域中间件, 让请求支持跨域-生产勿用
	// r.Use(middleware.Cors())
	// 初始化jwt auth中间件
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		panic(fmt.Sprintf("初始化jwt auth中间件失败: %v", err))
	}
	apiGroup := r.Group(common.Conf.System.UrlPathPrefix)
	// 注册公共路由，所有人都可以访问
	routers.InitPublicRouter(apiGroup)
	system.InitAuthRouter(apiGroup, authMiddleware) // 注册认证路由, 不会鉴权
	// 方便统一添加路由前缀
	v1 := apiGroup.Group("v1")
	{
		system.InitUserRouter(v1, authMiddleware)           // 注册用户路由
		system.InitDeptRouter(v1, authMiddleware)           // 注册部门路由
		system.InitMenuRouter(v1, authMiddleware)           // 注册菜单路由
		system.InitRoleRouter(v1, authMiddleware)           // 注册角色路由
		system.InitDictRouter(v1, authMiddleware)           // 注册字典路由
		system.InitAPIRouter(v1, authMiddleware)            // 接口api路由
		system.InitOperLogRouter(v1, authMiddleware)        // 注册操作日志路由
		cmdb.InitHostRouter(v1, authMiddleware)             // 注册主机路由
		prometheus.InitPrometheusRouter(v1, authMiddleware) // 注册prometheus路由
		k8s.InitPrometheusRouter(v1, authMiddleware)        // 注册k8sapi路由
	}

	return r
}
