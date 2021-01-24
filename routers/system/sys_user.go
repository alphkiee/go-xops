package system

import (
	"go-xops/api/v1/system"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// User
func InitUserRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/info", system.GetUserInfo)
		router.POST("/info/uploadImg", system.UserAvatarUpload)
		router.GET("/list", system.GetUsers)
		router.POST("/create", system.CreateUser)
		router.PATCH("/update/:userId", system.UpdateUserById)
		router.PATCH("/info/update/:userId", system.UpdateUserBaseInfoById)
		router.PUT("/changePwd", system.ChangePwd)
		router.DELETE("/delete", system.DeleteUserByIds)
	}
	return router
}
