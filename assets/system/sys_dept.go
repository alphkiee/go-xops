package system

import "go-xops/assets"

// 系统部门表
type SysDept struct {
	assets.Model
	Name     string    `gorm:"comment:'部门名称';size:64" json:"name"`
	Status   *bool     `gorm:"type:tinyint(1);default:1;comment:'状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Creator  string    `gorm:"comment:'创建人';size:64" json:"creator"`
	Sort     int       `gorm:"type:int(3);comment:'排序'" json:"sort"`
	ParentId uint      `gorm:"default:0;comment:'父级部门(编号为0时表示根)'" json:"parent_id"`
	Children []SysDept `gorm:"-" json:"children"` // 下属部门集合
	Users    []SysUser `gorm:"foreignkey:DeptId"` // 一个部门有多个user
}

func (m SysDept) TableName() string {
	return m.Model.TableName("sys_dept")
}
