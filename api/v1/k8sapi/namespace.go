package k8sapi

import (
	"context"
	"go-xops/dto/response"
	"go-xops/dto/service/k8ss"

	"github.com/gin-gonic/gin"
)

func GetNameSpaces(c *gin.Context) {
	result, err := k8ss.GetNameSpaces(context.TODO())
	if err != nil {
		response.FailWithMsg("获取namespace失败")
		return
	}
	response.SuccessWithData(result)
}

func GetNameSpace(c *gin.Context) {
	name := c.Param("ns")
	res, err := k8ss.GetNamespace(context.TODO(), name)
	if err != nil {
		response.FailWithMsg("获取namespace 失败")
		return
	}
	response.SuccessWithData(res)
}
