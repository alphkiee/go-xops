package k8s

import (
	"context"
	k8ss "go-xops/internal/service/k8s"
	"go-xops/pkg/common"

	"github.com/gin-gonic/gin"
)

// GetApplication doc
// @Summary Post /api/v1/k8s/application
// @Description 获取指定的application
// @Produce json
// @Param id query int false "id"
// @Param namespace query string false "namespace"
// @Param name query string false "name"
// @Param format query string false "format"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/k8s/application [post]
func GetApplication(c *gin.Context) {
	var req k8ss.ApplicationReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	res, err := k8ss.GetApplication(context.TODO(), req.NameSpace, req.ID, req.Name, req.Format)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.SuccessWithData(res)
}

// GetApplications doc
// @Summary Get /api/v1/k8s/application/list
// @Description 获取所有applications
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/k8s/application/list [get]
func GetApplications(c *gin.Context) {
	res, err := k8ss.GetApplications(context.TODO())
	if err != nil {
		common.FailWithMsg(err.Error())
	}
	common.SuccessWithData(res)
}
