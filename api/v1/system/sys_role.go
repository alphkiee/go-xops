package system

import (
	"fmt"
	"go-xops/assets/system"
	"go-xops/internal/request"
	"go-xops/internal/response"
	"go-xops/internal/service"
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
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/role/list [get]
func GetRoles(c *gin.Context) {
	// 绑定参数
	var req request.RoleListReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	roles, err := s.GetRoles(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.RoleListResp
	utils.Struct2StructByJson(roles, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response.SuccessWithData(resp)
}

// CreateRole doc
// @Summary Post /api/v1/role/create
// @Description 创建角色
// @Produce json
// @Param data body request.CreateRoleReq true "name, keyword, desc, creator"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/role/create [post]
func CreateRole(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateRoleReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	// 创建服务
	s := service.New()
	err = s.CreateRole(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateRoleById doc
// @Summary Patch /api/v1/role//update/:roleId
// @Description 根据role ID来更新角色基本信息
// @Produce json
// @Param roleId path int true "roleId"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/role/update/{roleId} [patch]
func UpdateRoleById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("角色编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateRoleById(roleId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateRolePermsById doc
// @Summary Patch /api/v1/role/perms/update/:roleId
// @Description 根据角色 ID来更新角色权限信息
// @Produce json
// @Param roleId path int true "roleId"
// @Param data body request.UpdateRolePermsReq true "menus_id, apis_id"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/role/perms/update/{roleId} [patch]
func UpdateRolePermsById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateRolePermsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("参数绑定失败, %v", err))
		return
	}
	// 获取path中的roleId
	roleId := utils.Str2Uint(c.Param("roleId"))
	if roleId == 0 {
		response.FailWithMsg("角色编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	if req.MenusId != nil {
		// 更新数据
		err = s.UpdateRoleMenusById(roleId, req.MenusId)
		if err != nil {
			response.FailWithMsg(err.Error())
			return
		}
	}
	if req.ApisId != nil {
		err = s.UpdateRoleApisById(roleId, req.ApisId)
		if err != nil {
			response.FailWithMsg(err.Error())
			return
		}
	}
	response.Success()
}

// BatchDeleteRoleByIds doc
// @Summary Delete /api/v1/role/delete
// @Description 根据ID批量删除角色
// @Produce json
// @Param data body request.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/role/delete [delete]
func BatchDeleteRoleByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteRoleByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// GetPermsByRoleId doc
// @Summary Get /api/v1/role/perms/:roleId
// @Description 获取当前用户信息
// @Produce json
// @Param roleId path int true "roleId"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/role/perms/:roleId [get]
func GetPermsByRoleId(c *gin.Context) {
	// 创建服务
	s := service.New()
	// 绑定参数
	resp, err := s.GetPermsByRoleId(utils.Str2Uint(c.Param("roleId")))
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.SuccessWithData(resp)
}
