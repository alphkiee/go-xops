package system

import (
	"go-xops/api/v1/system"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitRoleRouter ...角色路由
func InitRoleRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", system.GetRoles)
		router.POST("/create", system.CreateRole)
		router.PATCH("/update/:roleId", system.UpdateRoleById)
		router.PATCH("/perms/update/:roleId", system.UpdateRolePermsById)
		router.GET("/perms/:roleId", system.GetPermsByRoleId)
		router.DELETE("/delete", system.BatchDeleteRoleByIds)
	}
	return router
}
