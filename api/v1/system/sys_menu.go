package system

import (
	"go-xops/dto/request"
	"go-xops/dto/response"
	"go-xops/dto/service"
	"go-xops/models/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// 查询当前用户菜单树
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

// 查询所有菜单
func GetMenus(c *gin.Context) {
	// 创建服务
	s := service.New()
	menus := s.GetMenus()
	var resp []response.MenuTreeResp
	resp = service.GenMenuTree(nil, menus)
	response.SuccessWithData(resp)
}

// 创建菜单
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

// 更新菜单
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

// 批量删除菜单
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
