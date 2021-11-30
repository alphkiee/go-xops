package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	k8s "go-xops/internal/service/k8s"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// CreateDeplooyment doc
// @Summary Post /api/v1/k8s/deployment/create
// @Description 创建deployment
// @Produce json
// @Param deployment query string false "deployment"
// @Param replicas query int32 false "replicas"
// @Param matchlabels query string false "matchlabels"
// @Param labels query string false "labels"
// @Param imagename query string false "imagename"
// @Param image query string false "image"
// @Param namespace query string false "namespace"
// @Param containerport query int false "containerport"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/k8s/deployment/create [post]
func CreateDeplooyment(c *gin.Context) {
	var req k8s.DeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name: req.DeployName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(req.Replicas),
			Selector: &v1.LabelSelector{
				MatchLabels: req.MatchLabels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: req.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  req.ImageName,
							Image: req.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: int32(req.ContainerPort),
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := common.ClientSet.AppsV1().Deployments(req.NameSpace).Create(context.TODO(), deployment, v1.CreateOptions{})
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.SuccessWithMsg(data.GetObjectMeta().GetName())
}

// UpdateDeployment doc
// @Summary Patch /api/v1/k8s/deployment/update
// @Description 更新deployment
// @Produce json
// @Param deploymentname query string false "deploymentname"
// @Param replicas query int32 false "replicas"
// @Param image query string false "image"
// @Param namespace query string false "namespace"
// @Param containerport query int false "containerport"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/k8s/deployment/update [patch]
func UpdateDeployment(c *gin.Context) {
	var req k8s.UpDeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := common.ClientSet.AppsV1().Deployments(req.NameSpace).Get(context.TODO(), req.DeploymentName, v1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = utils.Int32Ptr(req.Replicas)       // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = req.Image // change nginx version
		_, updateErr := common.ClientSet.AppsV1().Deployments(req.NameSpace).Update(context.TODO(), result, v1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	common.SuccessWithMsg("滚动更新成功")
}

// DeleteDeployment doc
// @Summary Delete /api/v1/k8s/deployment/delete
// @Description 删除deployment
// @Produce json
// @Param namespace query string false "namespace"
// @Param name query string false "name"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/k8s/deployment/delete [delete]
func DeleteDeployment(c *gin.Context) {
	var req k8s.DeleteDeployment
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
	}
	// 删除deployment
	err = common.ClientSet.AppsV1().Deployments(req.Namespaces).Delete(context.TODO(), req.Name, v1.DeleteOptions{})
	if err != nil {
		common.FailWithMsg(err.Error())
	}
	common.Success()
}
