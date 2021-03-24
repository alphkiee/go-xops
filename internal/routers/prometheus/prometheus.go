package prometheus

import (
	"go-xops/api/v1/prometheus"
	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitPrometheusRouter ...
func InitPrometheusRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("prometheus").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/host/:key/:job", prometheus.Pmt)
	}
	return router
}
