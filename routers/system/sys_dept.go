package system

import (
	"go-xops/api/v1/system"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitDeptRouter ...部门路由
func InitDeptRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("dept").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/list", system.GetDepts)
		router.POST("/create", system.CreateDept)
		router.PATCH("/update/:deptId", system.UpdateDeptById)
		router.DELETE("/delete", system.BatchDeleteDeptByIds)
	}
	return router
}
