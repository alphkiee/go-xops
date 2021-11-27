package system

import (
	"go-xops/assets/system"
	s "go-xops/internal/service/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetUserMenuTree doc
// @Summary Get /api/v1/menu/tree
// @Description 当前用户菜单树
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/menu/tree [get]
func GetUserMenuTree(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	menus, err := s.GetUserMenuList(user.(system.SysUser).RoleId)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	var resp []s.MenuTreeResp
	// 转换成树结构
	resp = s.GenMenuTree(nil, menus)
	common.SuccessWithData(resp)
}

// GetMenus doc
// @Summary Get /api/v1/menu/list
// @Description 查询所有菜单
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Router /api/v1/menu/list [get]
func GetMenus(c *gin.Context) {
	menus := s.GetMenus()
	var resp []s.MenuTreeResp
	resp = s.GenMenuTree(nil, menus)
	common.SuccessWithData(resp)
}

// CreateMenu doc
// @Summary Get /api/v1/menu/create
// @Description 创建菜单
// @Produce json
// @Param data body s.CreateMenuReq true "name, icon, path, sort, parent_id, creator"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/menu/create [post]
func CreateMenu(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req s.CreateMenuReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["Name"] = "姓名"

	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	req.Creator = user.(system.SysUser).Name
	err = s.CreateMenu(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateMenuById doc
// @Summary Get /api/v1/menu/update/:menuId
// @Description 更新菜单
// @Produce json
// @Param data body s.UpdateMenuReq true "name, icon, path, sort, status, parent_id"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/menu/update/:menuId [patch]
func UpdateMenuById(c *gin.Context) {
	// 绑定参数
	var req s.UpdateMenuReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}

	// 获取path中的menuId
	menuId := utils.Str2Uint(c.Param("menuId"))
	if menuId == 0 {
		common.FailWithMsg("菜单编号不正确")
		return
	}
	err = s.UpdateMenuById(menuId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// BatchDeleteMenuByIds doc
// @Summary Delete /api/v1/menu/delete
// @Description 根据ID批量删除菜单
// @Produce json
// @Param data body common.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/menu/delete [delete]
func BatchDeleteMenuByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = s.DeleteMenuByIds(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}
