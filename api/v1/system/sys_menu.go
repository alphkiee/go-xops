package system

import (
	"go-xops/assets/system"
	"go-xops/internal/request"
	"go-xops/internal/response"
	"go-xops/internal/service"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetUserMenuTree doc
// @Summary Get /api/v1/menu/tree
// @Description 当前用户菜单树
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/menu/tree [get]
func GetUserMenuTree(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	s := service.New()
	menus, err := s.GetUserMenuList(user.(system.SysUser).RoleId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	var resp []response.MenuTreeResp
	// 转换成树结构
	resp = service.GenMenuTree(nil, menus)
	response.SuccessWithData(resp)
}

// GetMenus doc
// @Summary Get /api/v1/menu/list
// @Description 查询所有菜单
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Router /api/v1/menu/list [get]
func GetMenus(c *gin.Context) {
	// 创建服务
	s := service.New()
	menus := s.GetMenus()
	var resp []response.MenuTreeResp
	resp = service.GenMenuTree(nil, menus)
	response.SuccessWithData(resp)
}

// CreateMenu doc
// @Summary Get /api/v1/menu/create
// @Description 创建菜单
// @Produce json
// @Param data body request.CreateMenuReq true "name, icon, path, sort, parent_id, creator"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/menu/create [post]
func CreateMenu(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateMenuReq
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
	err = s.CreateMenu(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateMenuById doc
// @Summary Get /api/v1/menu/update/:menuId
// @Description 更新菜单
// @Produce json
// @Param data body request.UpdateMenuReq true "name, icon, path, sort, status, parent_id"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/menu/update/:menuId [patch]
func UpdateMenuById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateMenuReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 获取path中的menuId
	menuId := utils.Str2Uint(c.Param("menuId"))
	if menuId == 0 {
		response.FailWithMsg("菜单编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateMenuById(menuId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteMenuByIds doc
// @Summary Delete /api/v1/menu/delete
// @Description 根据ID批量删除菜单
// @Produce json
// @Param data body request.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/menu/delete [delete]
func BatchDeleteMenuByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteMenuByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
