package system

import "go-xops/models"

// 系统角色表
type SysRole struct {
	models.Model
	Name    string    `gorm:"comment:'角色名称';size:128" json:"name"`
	Keyword string    `gorm:"uniqueIndex:uk_keyword;comment:'角色关键词';size:64" json:"keyword"`
	Desc    string    `gorm:"comment:'角色说明';size:255" json:"desc"`
	Status  *bool     `gorm:"type:tinyint(1);default:1;comment:'角色状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Creator string    `gorm:"comment:'创建人';size:128" json:"creator"`
	Menus   []SysMenu `gorm:"many2many:relation_role_menu;" json:"menus"` // 角色菜单多对多关系
	Users   []SysUser `gorm:"foreignkey:RoleId"`                          // 一个角色有多个user
	//Apis    []SysApi  `gorm:"many2many:relation_role_api;" json:"apis"`   // 角色接口多对多关系
}

func (m SysRole) TableName() string {
	return m.Model.TableName("sys_role")
}
