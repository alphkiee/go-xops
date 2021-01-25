package service

import (
	"errors"
	"fmt"
	"go-xops/dto/request"
	"go-xops/dto/response"
	"go-xops/models/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取所有角色
func (s *MysqlService) GetRoles(req *request.RoleListReq) ([]system.SysRole, error) {
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
func (s *MysqlService) GetPermsByRoleId(roleId uint) (response.RolePermsResp, error) {
	var role system.SysRole
	//var apis []models.SysApi
	var resp response.RolePermsResp
	err := s.db.Preload("Menus", "status = ?", true).Where("id = ?", roleId).First(&role).Error
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
	err = s.db.Find(&allApi).Error
	if err == nil {
		casbins, err := s.GetCasbinListByRoleId(roleId)
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
func (s *MysqlService) CreateRole(req *request.CreateRoleReq) (err error) {
	var role system.SysRole
	utils.Struct2StructByJson(req, &role)
	// 创建数据
	err = s.db.Create(&role).Error
	return
}

// 更新角色
func (s *MysqlService) UpdateRoleById(id uint, req gin.H) (err error) {
	var oldRole system.SysRole
	query := s.db.Model(oldRole).Where("id = ?", id).First(&oldRole)
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
func (s *MysqlService) UpdateRoleMenusById(id uint, req []uint) (err error) {
	var menus []system.SysMenu
	err = s.db.Where("id in (?)", req).Find(&menus).Error
	if err != nil {
		return
	}
	// 替换菜单
	var role system.SysRole
	err = s.db.Where("id = ?", id).First(&role).Error
	err = s.db.Model(&role).Association("Menus").Replace(&menus)
	return
}

// 更新角色的权限接口
func (s *MysqlService) UpdateRoleApisById(id uint, req []uint) (err error) {
	var oldRole system.SysRole
	query := s.db.Model(&oldRole).Where("id = ?", id).First(&oldRole)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("角色不存在")
	}
	if len(req) > 0 {
		// 先删除当前角色的规则
		oldCasbin := s.GetRoleCasbins(system.SysRoleCasbin{
			Keyword: oldRole.Keyword,
		})
		// 批量删除
		_, err = s.BatchDeleteRoleCasbins(oldCasbin)

		// 构建新的规则
		apis := make([]system.SysApi, 0)
		err = s.db.Where("id IN (?)", req).Find(&apis).Error
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
		_, err = s.BatchCreateRoleCasbins(cs)
	}
	return
}

// 批量删除角色
func (s *MysqlService) DeleteRoleByIds(ids []uint) (err error) {
	var roles []system.SysRole
	// 查询符合条件的角色, 以及关联的用户
	err = s.db.Preload("Users").Preload("Menus").Where("id IN (?)", ids).Find(&roles).Error
	if err != nil {
		return
	}
	newIds := make([]uint, 0)
	oldCasbins := make([]system.SysRoleCasbin, 0)
	for _, role := range roles {
		if len(role.Users) > 0 {
			return errors.New(fmt.Sprintf("角色[%s]仍有%d位关联用户, 请先移除关联用户再删除角色", role.Name, len(role.Users)))
		}
		oldCasbins = append(oldCasbins, s.GetRoleCasbins(system.SysRoleCasbin{
			Keyword: role.Keyword,
		})...)
		newIds = append(newIds, role.Id)
	}
	if len(oldCasbins) > 0 {
		// 删除关联的casbin
		_, err = s.BatchDeleteRoleCasbins(oldCasbins)
	}
	if len(newIds) > 0 {
		// 执行删除
		err = s.db.Where("id IN (?)", newIds).Delete(system.SysRole{}).Error
	}
	return
}
