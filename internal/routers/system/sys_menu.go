package system

import (
	"go-xops/api/v1/system"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitMenuRouter ...菜单路由
func InitMenuRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("menu").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/tree", system.GetUserMenuTree)
		router.GET("/list", system.GetMenus)
		router.POST("/create", system.CreateMenu)
		router.PATCH("/update/:menuId", system.UpdateMenuById)
		router.DELETE("/delete", system.BatchDeleteMenuByIds)
	}
	return router
}
