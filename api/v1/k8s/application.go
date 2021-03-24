package k8s

import (
	"context"
	"go-xops/internal/request/k8s"
	"go-xops/internal/response"
	k8ss "go-xops/internal/service/k8s"

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
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/k8s/application [post]
func GetApplication(c *gin.Context) {
	var req k8s.ApplicationReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	res, err := k8ss.GetApplication(context.TODO(), req.NameSpace, req.ID, req.Name, req.Format)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(res)
}

// GetApplications doc
// @Summary Get /api/v1/k8s/application/list
// @Description 获取所有applications
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/k8s/application/list [get]
func GetApplications(c *gin.Context) {
	res, err := k8ss.GetApplications(context.TODO())
	if err != nil {
		response.FailWithMsg(err.Error())
	}
	response.SuccessWithData(res)
}
