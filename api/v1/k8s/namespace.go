package k8s

import (
	"context"
	"go-xops/internal/response"
	k8ss "go-xops/internal/service/k8s"

	"github.com/gin-gonic/gin"
)

// GetNameSpaces doc
// @Summary Get /api/v1/k8s/namespace/list
// @Description 获取所有的namespaces
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/k8s/namespace/list [get]
func GetNameSpaces(c *gin.Context) {
	result, err := k8ss.GetNameSpaces(context.TODO())
	if err != nil {
		response.FailWithMsg("获取namespace失败")
		return
	}
	response.SuccessWithData(result)
}

// GetNameSpace doc
// @Summary Get /api/v1/k8s/namespace/name/:ns
// @Description 获取指定的namespace
// @Produce json
// @Param ns path string ture "ns"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/k8s/namespace/name/:ns [get]
func GetNameSpace(c *gin.Context) {
	name := c.Param("ns")
	res, err := k8ss.GetNamespace(context.TODO(), name)
	if err != nil {
		response.FailWithMsg("获取namespace 失败")
		return
	}
	response.SuccessWithData(res)
}
