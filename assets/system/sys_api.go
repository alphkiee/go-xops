package system

import "go-xops/assets"

// 系统接口表
type SysApi struct {
	assets.Model
	Name     string `gorm:"comment:'接口名称';size:64" json:"name"`
	Method   string `gorm:"comment:'请求方式';size:64" json:"method"`
	Path     string `gorm:"comment:'访问路径';size:128" json:"path"`
	Category string `gorm:"comment:'所属类别';size:128" json:"category"`
	Desc     string `gorm:"comment:'说明';size:128" json:"desc"`
	Creator  string `gorm:"comment:'创建人';size:64" json:"creator"`
	//Roles      []SysRole `gorm:"many2many:relation_role_api;" json:"roles"` // 角色接口多对多关系
}

func (m SysApi) TableName() string {
	return m.Model.TableName("sys_api")
}
