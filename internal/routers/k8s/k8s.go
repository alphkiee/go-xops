package k8s

import (
	"go-xops/api/v1/k8s"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitPrometheusRouter ...
func InitPrometheusRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("k8s").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/namespace/list", k8s.GetNameSpaces)
		router.GET("/namespace/name/:ns", k8s.GetNameSpace)
		router.POST("/application", k8s.GetApplication)
		router.GET("/application/list", k8s.GetApplications)
		router.POST("/deployment/create", k8s.CreateDeplooyment)
		router.PATCH("/deployment/update", k8s.UpdateDeployment)
		router.DELETE("/deployment/delete", k8s.DeleteDeployment)
		router.POST("/pod/exec", k8s.ExecPod)
	}
	return router
}
