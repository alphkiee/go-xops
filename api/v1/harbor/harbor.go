package harbor

import (
	"github.com/gin-gonic/gin"
	ha "go-xops/assets/harbor"
	"go-xops/internal/service/harbor"
	"go-xops/pkg/common"
)

// GetProjects doc
// @Summary Get /api/v1/harbor/projects/list
// @Description 获取所有projects
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/harbor/projects/list [get]
func GetProjects(c *gin.Context) {
	/*
		var rep harbor.ProjectsListReq
		err := c.ShouldBindJSON(rep)
		if err != nil {
			common.FailWithMsg(err.Error())
			return
		}
	*/
	p, err := harbor.GetProjectList()
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.SuccessWithData(p)
}

// CreateProject doc
// @Summary Post /api/v1/harbor/project/create
// @Description 创建project
// @Produce json
// @Param name query string false "name"
// @Param storage_limit query string false "storage_limit"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/harbor/project/create [post]
func CreateProject(c *gin.Context) {
	var req harbor.ProjectCreateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	P, err := harbor.CreateProject(req.Name, req.CountLimit, req.StorageLimit)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.SuccessWithData(P)
}

// DeleteProject doc
// @Summary Delete /api/v1/harbor/project/delete
// @Description 删除project
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/harbor/project/delete [delete]
func DeleteProject(c *gin.Context) {
	var req harbor.ProjectDeleteReq
	var p ha.Project
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	p.Name = req.Name
	p.ProjectID = req.ID
	err = harbor.DeleteProject(p)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateProject doc
// @Summary Put /api/v1/harbor/project/update
// @Description 更新project
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/harbor/project/update [put]
func UpdateProject(c *gin.Context) {
	var req harbor.ProjectUpdate
	var p ha.Project
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	p.UpdateTime = req.UpdateTime
	p.ProjectID = req.ProjectID
	p.OwnerID = req.OwnerID
	p.CreationTime = req.CreationTime
	p.RepoCount = req.RepoCount
	p.Name = req.Name
	p.CurrentUserRoleID = req.CurrentUserRoleID
	err = harbor.UpdateProject(p, req.CountLimit, req.StorageLimit)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// GetRegistry doc
// @Summary Get /api/v1/harbor/registry/:name
// @Description 获取所有Registry
// @Produce json
// @Param name path string true "name"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/harbor/registry/:name [get]
func GetRegistry(c *gin.Context) {
	P, err := harbor.GetRegistry(c.Param("name"))
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.SuccessWithData(P)
}

// GetRepository doc
// @Summary Get /api/v1/harbor/repository/:name
// @Description 获取指定project目录下的镜像
// @Produce json
// @Param name path string true "name"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/harbor/repository/:name [get]
func GetRepository(c *gin.Context) {
	R, err := harbor.GetSearch(c.Param("name"))
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}

	common.SuccessWithData(R)
}
