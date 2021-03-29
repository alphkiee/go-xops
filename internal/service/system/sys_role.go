package system

import (
	"errors"
	"fmt"
	"go-xops/assets"
	"go-xops/assets/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取角色列表结构体
type RoleListReq struct {
	Name            string `json:"name" form:"name"`
	Keyword         string `json:"keyword" form:"keyword"`
	Status          *bool  `json:"status" form:"status"`
	Creator         string `json:"creator" form:"creator"`
	common.PageInfo        // 分页参数
}

// 更新角色权限的结构体
type UpdateRolePermsReq struct {
	MenusId []uint `json:"menus_id" form:"menus_id"` // 传多个id
	ApisId  []uint `json:"apis_id" form:"apis_id"`   // 传多个id
}

// 创建角色结构体
type CreateRoleReq struct {
	Name    string `json:"name" validate:"required"`
	Keyword string `json:"keyword" validate:"required"`
	Desc    string `json:"desc"`
	Creator string `json:"creator"`
}

// 角色返回权限信息
type RolePermsResp struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
	MenusId []uint `json:"menus_id"`
	ApisId  []uint `json:"apis_id"`
}

// 角色信息响应, 字段含义见assets
type RoleListResp struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	Keyword   string           `json:"keyword"`
	Desc      string           `json:"desc"`
	Status    *bool            `json:"status"`
	Creator   string           `json:"creator"`
	CreatedAt assets.LocalTime `json:"created_at"`
}

// 获取所有角色
func GetRoles(req *RoleListReq) ([]system.SysRole, error) {
	var err error
	list := make([]system.SysRole, 0)
	db := common.Mysql
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	// 查询条数
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.All {
			// 不使用分页
			err = db.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = db.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 根据角色ID获取角色权限：菜单和接口
func GetPermsByRoleId(roleId uint) (RolePermsResp, error) {
	var role system.SysRole
	var resp RolePermsResp
	err := common.Mysql.Preload("Menus", "status = ?", true).Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return resp, err
	}
	resp.Id = role.Id
	resp.Name = role.Name
	resp.Keyword = role.Keyword
	for _, menu := range role.Menus {
		resp.MenusId = append(resp.MenusId, menu.Id)
	}
	allApi := make([]system.SysApi, 0)
	// 查询全部api
	err = common.Mysql.Find(&allApi).Error
	if err == nil {
		casbins, err := GetCasbinListByRoleId(roleId)
		if err == nil {
			for _, api := range allApi {
				path := api.Path
				method := api.Method
				for _, casbin := range casbins {
					// 该api有权限
					if path == casbin.V1 && method == casbin.V2 {
						resp.ApisId = append(resp.ApisId, api.Id)
						break
					}
				}
			}
		}
	}
	return resp, err
}

// 创建角色
func CreateRole(req *CreateRoleReq) (err error) {
	var role system.SysRole
	utils.Struct2StructByJson(req, &role)
	// 创建数据
	err = common.Mysql.Create(&role).Error
	return
}

// 更新角色
func UpdateRoleById(id uint, req gin.H) (err error) {
	var oldRole system.SysRole
	query := common.Mysql.Model(oldRole).Where("id = ?", id).First(&oldRole)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	// 比对增量字段
	var m system.SysRole
	utils.CompareDifferenceStructByJson(oldRole, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 更新角色的菜单
func UpdateRoleMenusById(id uint, req []uint) (err error) {
	var menus []system.SysMenu
	err = common.Mysql.Where("id in (?)", req).Find(&menus).Error
	if err != nil {
		return
	}
	// 替换菜单
	var role system.SysRole
	err = common.Mysql.Where("id = ?", id).First(&role).Error
	err = common.Mysql.Model(&role).Association("Menus").Replace(&menus)
	return
}

// 更新角色的权限接口
func UpdateRoleApisById(id uint, req []uint) (err error) {
	var oldRole system.SysRole
	query := common.Mysql.Model(&oldRole).Where("id = ?", id).First(&oldRole)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("角色不存在")
	}
	if len(req) > 0 {
		// 先删除当前角色的规则
		oldCasbin := GetRoleCasbins(system.SysRoleCasbin{
			Keyword: oldRole.Keyword,
		})
		// 批量删除
		_, err = BatchDeleteRoleCasbins(oldCasbin)

		// 构建新的规则
		apis := make([]system.SysApi, 0)
		err = common.Mysql.Where("id IN (?)", req).Find(&apis).Error
		if err != nil {
			return
		}
		cs := make([]system.SysRoleCasbin, 0)
		for _, api := range apis {
			cs = append(cs, system.SysRoleCasbin{
				Keyword: oldRole.Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})
		}
		// 批量创建规则
		_, err = BatchCreateRoleCasbins(cs)
	}
	return
}

// 批量删除角色
func DeleteRoleByIds(ids []uint) (err error) {
	var roles []system.SysRole
	// 查询符合条件的角色, 以及关联的用户
	err = common.Mysql.Preload("Users").Preload("Menus").Where("id IN (?)", ids).Find(&roles).Error
	if err != nil {
		return
	}
	newIds := make([]uint, 0)
	oldCasbins := make([]system.SysRoleCasbin, 0)
	for _, role := range roles {
		if len(role.Users) > 0 {
			return errors.New(fmt.Sprintf("角色[%s]仍有%d位关联用户, 请先移除关联用户再删除角色", role.Name, len(role.Users)))
		}
		oldCasbins = append(oldCasbins, GetRoleCasbins(system.SysRoleCasbin{
			Keyword: role.Keyword,
		})...)
		newIds = append(newIds, role.Id)
	}
	if len(oldCasbins) > 0 {
		// 删除关联的casbin
		_, err = BatchDeleteRoleCasbins(oldCasbins)
	}
	if len(newIds) > 0 {
		// 执行删除
		err = common.Mysql.Where("id IN (?)", newIds).Delete(system.SysRole{}).Error
	}
	return
}
