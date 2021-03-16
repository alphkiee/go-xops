package k8sr

import (
	"go-xops/api/v1/k8sapi"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitPrometheusRouter ...
func InitPrometheusRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("k8s").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/namespace/list", k8sapi.GetNameSpaces)
		router.GET("/namespace/name/:ns", k8sapi.GetNameSpace)
		router.POST("/application", k8sapi.GetApplication)
		router.GET("/application/list", k8sapi.GetApplications)
		router.POST("/deployment/create", k8sapi.CreateDeplooyment)
		router.PATCH("/deployment/update", k8sapi.UpdateDeployment)
	}
	return router
}
