package system

import "go-xops/assets"

// SysDict ...
type SysDict struct {
	assets.Model
	Key      string    `gorm:"uniqueIndex:uk_key;comment:'字典Key';size:64" json:"key"`
	Value    string    `gorm:"comment:'字典Value';size:64" json:"value"`
	Desc     string    `gorm:"comment:'说明';size:128" json:"desc"`
	Status   *bool     `gorm:"type:tinyint(1);default:1;comment:'状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Creator  string    `gorm:"comment:'创建人';size:64" json:"creator"`
	ParentId uint      `gorm:"default:0;comment:'父级字典(编号为0时表示根)'" json:"parent_id"`
	Dicts    []SysDict `gorm:"foreignkey:ParentId"`
}

func (m SysDict) TableName() string {
	return m.Model.TableName("sys_dict")
}
