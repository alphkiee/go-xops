package system

import (
	"fmt"
	"go-xops/assets/system"
	s "go-xops/internal/service/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetRoles doc
// @Summary Get /api/v1/role/list
// @Description 列出所有角色
// @Produce json
// @Param name query string false "name"
// @Param keyword query string false "keyword"
// @Param status query string false "status"
// @Param creator query string false "creator"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/role/list [get]
func GetRoles(c *gin.Context) {
	// 绑定参数
	var req s.RoleListReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	roles, err := s.GetRoles(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []s.RoleListResp
	utils.Struct2StructByJson(roles, &respStruct)
	// 返回分页数据
	var resp common.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	common.SuccessWithData(resp)
}

// CreateRole doc
// @Summary Post /api/v1/role/create
// @Description 创建角色
// @Produce json
// @Param data body s.CreateRoleReq true "name, keyword, desc, creator"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/role/create [post]
func CreateRole(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req s.CreateRoleReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["Name"] = "姓名"
	m["Keyword"] = "关键字"
	m["Desc"] = "描述"
	m["Creator"] = "创建人"
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	req.Creator = user.(system.SysUser).Name
	err = s.CreateRole(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateRoleById doc
// @Summary Patch /api/v1/role//update/:roleId
// @Description 根据role ID来更新角色基本信息
// @Produce json
// @Param roleId path int true "roleId"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/role/update/{roleId} [patch]
func UpdateRoleById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		common.FailWithMsg("角色编号不正确")
		return
	}
	err = s.UpdateRoleById(roleId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateRolePermsById doc
// @Summary Patch /api/v1/role/perms/update/:roleId
// @Description 根据角色 ID来更新角色权限信息
// @Produce json
// @Param roleId path int true "roleId"
// @Param data body s.UpdateRolePermsReq true "menus_id, apis_id"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/role/perms/update/{roleId} [patch]
func UpdateRolePermsById(c *gin.Context) {
	// 绑定参数
	var req s.UpdateRolePermsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithMsg(fmt.Sprintf("参数绑定失败, %v", err))
		return
	}
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		common.FailWithMsg("角色编号不正确")
		return
	}
	if req.MenusId != nil {
		// 更新数据
		err = s.UpdateRoleMenusById(roleId, req.MenusId)
		if err != nil {
			common.FailWithMsg(err.Error())
			return
		}
	}
	if req.ApisId != nil {
		err = s.UpdateRoleApisById(roleId, req.ApisId)
		if err != nil {
			common.FailWithMsg(err.Error())
			return
		}
	}
	common.Success()
}

// BatchDeleteRoleByIds doc
// @Summary Delete /api/v1/role/delete
// @Description 根据ID批量删除角色
// @Produce json
// @Param data body common.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/role/delete [delete]
func BatchDeleteRoleByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = s.DeleteRoleByIds(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// GetPermsByRoleId doc
// @Summary Get /api/v1/role/perms/:roleId
// @Description 获取当前用户信息
// @Produce json
// @Param roleId path int true "roleId"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/role/perms/:roleId [get]
func GetPermsByRoleId(c *gin.Context) {
	resp, err := s.GetPermsByRoleId(utils.Str2Uint(c.Param("roleId")))
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.SuccessWithData(resp)
}
