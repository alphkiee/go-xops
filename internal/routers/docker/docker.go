package docker

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-xops/api/v1/docker"
	"go-xops/middleware"
)

// docker reset api接口调用路由
func InitDockerRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("docker").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		//router.GET("/build", docker.BuildImage)
		router.GET("/build/socket", docker.BuildImageSocket)
		router.POST("/push/image", docker.PushImage)
		router.GET("push/image/socket", docker.BuildImageSocketPush)
	}
	return router
}
