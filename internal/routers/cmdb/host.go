package cmdb

import (
	"go-xops/api/v1/cmdb"

	"go-xops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitHostRouter ...
func InitHostRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("host").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.GET("/info/:id", cmdb.GetHostInfo)
		router.GET("/list", cmdb.GetHosts)
		router.POST("/create", cmdb.CreateHost)
		router.PATCH("update/:id", cmdb.UpdateHostById)
		router.DELETE("/delete", cmdb.BatchDeleteHostByIds)
		router.GET("/ssh/:sid", cmdb.WsSsh)
		router.GET("/sftp/:sid", cmdb.Sftp_ssh)
		router.POST("/upload", cmdb.ExcelIn)
		router.GET("/exphost/:ids", cmdb.ExportHost)
		router.GET("/cmd/:ids/:cmds", cmdb.CmdExec)
		router.POST("/fileUpload", cmdb.Sftp)
		router.POST("/fileDownload", cmdb.SftpDow)
	}
	return router
}
