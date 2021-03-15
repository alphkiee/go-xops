package k8sapi

import (
	"context"
	"go-xops/dto/request/k8s_request"
	"go-xops/dto/response"
	"go-xops/dto/service/k8ss"

	"github.com/gin-gonic/gin"
)

func GetApplication(c *gin.Context) {
	var req k8s_request.ApplicationReq
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

func GetApplications(c *gin.Context) {
	res, err := k8ss.GetApplications(context.TODO())
	if err != nil {
		response.FailWithMsg(err.Error())
	}
	response.SuccessWithData(res)
}
