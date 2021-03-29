package k8s

import (
	"context"
	k8ss "go-xops/internal/service/k8s"
	"go-xops/pkg/common"

	"github.com/gin-gonic/gin"
)

// GetNameSpaces doc
// @Summary Get /api/v1/k8s/namespace/list
// @Description 获取所有的namespaces
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/k8s/namespace/list [get]
func GetNameSpaces(c *gin.Context) {
	result, err := k8ss.GetNameSpaces(context.TODO())
	if err != nil {
		common.FailWithMsg("获取namespace失败")
		return
	}
	common.SuccessWithData(result)
}

// GetNameSpace doc
// @Summary Get /api/v1/k8s/namespace/name/:ns
// @Description 获取指定的namespace
// @Produce json
// @Param ns path string ture "ns"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/k8s/namespace/name/:ns [get]
func GetNameSpace(c *gin.Context) {
	name := c.Param("ns")
	res, err := k8ss.GetNamespace(context.TODO(), name)
	if err != nil {
		common.FailWithMsg("获取namespace 失败")
		return
	}
	common.SuccessWithData(res)
}
