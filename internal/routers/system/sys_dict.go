package system

import (
	"go-xops/api/v1/system"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitDictRouter ...字典路由
func InitDictRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("dict").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", system.GetDicts)
		router.POST("/create", system.CreateDict)
		router.PATCH("/update/:dictId", system.UpdateDictById)
		router.DELETE("/delete", system.BatchDeleteDictByIds)
	}
	return router
}
