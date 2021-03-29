package system

import "go-xops/assets"

// SysRole ...
type SysRole struct {
	assets.Model
	Name    string    `gorm:"comment:'角色名称';size:128" json:"name"`
	Keyword string    `gorm:"uniqueIndex:uk_keyword;comment:'角色关键词';size:64" json:"keyword"`
	Desc    string    `gorm:"comment:'角色说明';size:255" json:"desc"`
	Status  *bool     `gorm:"type:tinyint(1);default:1;comment:'角色状态(正常/禁用, 默认正常)'" json:"status"`
	Creator string    `gorm:"comment:'创建人';size:128" json:"creator"`
	Menus   []SysMenu `gorm:"many2many:relation_role_menu;" json:"menus"`
	Users   []SysUser `gorm:"foreignkey:RoleId"`
}

func (m SysRole) TableName() string {
	return m.Model.TableName("sys_role")
}
