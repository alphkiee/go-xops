package harbor

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-xops/api/v1/harbor"
	"go-xops/middleware"
)

// harbor reset api接口调用路由
func InitHarborRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("harbor").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/projects/list", harbor.GetProjects)
		router.POST("/project/create", harbor.CreateProject)
		router.DELETE("project/delete", harbor.DeleteProject)
		router.PUT("/project/update", harbor.UpdateProject)
		router.GET("/registry/:name", harbor.GetRegistry)
		router.GET("/repository/:name", harbor.GetRepository)
	}
	return router
}
