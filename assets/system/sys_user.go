package system

import "go-xops/assets"

// SysUser ...
type SysUser struct {
	assets.Model
	Username string  `gorm:"uniqueIndex:uk_username;comment:'用户名';size:128" json:"username" binding:"required"`
	Password string  `gorm:"comment:'密码';size:128" json:"password" binding:"required"`
	Mobile   string  `gorm:"comment:'手机';size:128" json:"mobile"`
	Avatar   string  `gorm:"comment:'头像';size:128" json:"avatar"`
	Name     string  `gorm:"comment:'姓名';size:128" json:"name"`
	Email    string  `gorm:"comment:'邮箱地址';size:128" json:"email"`
	Status   *bool   `gorm:"type:tinyint(1);default:1;comment:'用户状态(正常/禁用, 默认正常)'" json:"status"`
	Creator  string  `gorm:"comment:'创建人';size:128" json:"creator"`
	RoleId   uint    `gorm:"comment:'角色Id外键'" json:"role_id"`
	Role     SysRole `gorm:"foreignkey:RoleId" json:"role"`
	DeptId   uint    `gorm:"comment:'部门Id外键'" json:"dept_id"`
	Dept     SysDept `gorm:"foreignkey:DeptId" json:"dept"`
}

func (m SysUser) TableName() string {
	return m.Model.TableName("sys_user")
}
