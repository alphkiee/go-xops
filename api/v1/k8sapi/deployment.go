package k8sapi

import (
	"context"
	"fmt"
	"go-xops/dto/request/k8s_request"
	"go-xops/dto/response"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func CreateDeplooyment(c *gin.Context) {
	var req k8s_request.DeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
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
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithMsg(data.GetObjectMeta().GetName())
}

func UpdateDeployment(c *gin.Context) {
	var req k8s_request.UpDeploymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := common.ClientSet.AppsV1().Deployments(req.NameSpace).Get(context.TODO(), req.DeploymentName, v1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = utils.Int32Ptr(req.Replicas)          // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := common.ClientSet.AppsV1().Deployments(req.NameSpace).Update(context.TODO(), result, v1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	response.SuccessWithMsg("滚动更新成功")
}
