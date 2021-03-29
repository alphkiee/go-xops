package system

import "go-xops/assets"

// 系统菜单表
type SysMenu struct {
	assets.Model
	Name     string    `gorm:"comment:'菜单名称';size:64" json:"name"`
	Icon     string    `gorm:"comment:'菜单图标';size:64" json:"icon"`
	Path     string    `gorm:"comment:'菜单访问路径';size:64" json:"path"`
	Sort     int       `gorm:"type:int(3);comment:'菜单顺序(同级菜单, 从0开始, 越小显示越靠前)'" json:"sort"`
	Status   *bool     `gorm:"type:tinyint(1);default:1;comment:'菜单状态(正常/禁用, 默认正常)'" json:"status"`
	ParentId uint      `gorm:"default:0;comment:'父菜单编号(编号为0时表示根菜单)'" json:"parent_id"`
	Creator  string    `gorm:"comment:'创建人';size:64" json:"creator"`
	Children []SysMenu `gorm:"-" json:"children"`
	Roles    []SysRole `gorm:"many2many:relation_role_menu;" json:"roles"`
}

func (m SysMenu) TableName() string {
	return m.Model.TableName("sys_menu")
}

// 获取选中列表
func GetCheckedMenuIds(list []uint, allMenu []SysMenu) []uint {
	checked := make([]uint, 0)
	for _, c := range list {
		parent := SysMenu{
			ParentId: c,
		}
		children := parent.GetChildrenIds(allMenu)
		count := 0
		for _, child := range children {
			contains := false
			for _, v := range list {
				if v == child {
					contains = true
				}
			}
			if contains {
				count++
			}
		}
		if len(children) == count {
			checked = append(checked, c)
		}
	}
	return checked
}

// 查找子菜单编号
func (m SysMenu) GetChildrenIds(allMenu []SysMenu) []uint {
	childrenIds := make([]uint, 0)
	for _, menu := range allMenu {
		if menu.ParentId == m.ParentId {
			childrenIds = append(childrenIds, menu.Id)
		}
	}
	return childrenIds
}
