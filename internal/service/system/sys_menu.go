package system

import (
	"errors"
	"go-xops/assets/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"sort"

	"gorm.io/gorm"
)

type MenuListReq struct {
	Title           string `json:"title" form:"title"`
	Path            string `json:"path" form:"path"`
	Status          *bool  `json:"status" form:"status"`
	Creator         string `json:"creator" form:"creator"`
	common.PageInfo        // 分页参数
}

// 创建菜单结构体
type CreateMenuReq struct {
	Name     string `json:"name" validate:"required"`
	Icon     string `json:"icon"`
	Path     string `json:"path"`
	Sort     int    `json:"sort"`
	ParentId uint   `json:"parent_id"`
	Creator  string `json:"creator"`
}

// 修改菜单
type UpdateMenuReq struct {
	Name     string `json:"name" validate:"required"`
	Icon     string `json:"icon"`
	Path     string `json:"path"`
	Sort     int    `json:"sort"`
	Status   *bool  `json:"status"`
	ParentId uint   `json:"parent_id"`
}

// 菜单树信息响应, 字段含义见models.SysMenu
type MenuTreeResp struct {
	Id       uint           `json:"id"`
	ParentId uint           `json:"parent_id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon"`
	Path     string         `json:"path"`
	Creator  string         `json:"creator"`
	Sort     int            `json:"sort"`
	Status   bool           `json:"status"`
	Children []MenuTreeResp `json:"children,omitempty"` //tag:omitempty 为空的值不显示
}
type MenuTreeRespList []MenuTreeResp

func (hs MenuTreeRespList) Len() int {
	return len(hs)
}
func (hs MenuTreeRespList) Less(i, j int) bool {
	return hs[i].Sort < hs[j].Sort // 按Sort从小到大排序
}

func (hs MenuTreeRespList) Swap(i, j int) {
	hs[i], hs[j] = hs[j], hs[i]
}

// 获取用户菜单的切片
func GetUserMenuList(roleId uint) ([]system.SysMenu, error) {
	//tree := make([]models.SysMenu, 0)
	var role system.SysRole
	err := common.Mysql.Table(new(system.SysRole).TableName()).Preload("Menus", "status = ?", true).Where("id = ?", roleId).Find(&role).Error
	menus := make([]system.SysMenu, 0)
	if err != nil {
		return menus, err
	}
	if role.Keyword == "admin" {
		err = common.Mysql.Find(&menus).Error
		return menus, err
	}
	// 生成菜单树
	//tree = GenMenuTree(nil, role.Menus)
	return role.Menus, nil
}

// 获取所有菜单
func GetMenus() []system.SysMenu {
	//tree := make([]models.SysMenu, 0)
	menus := getAllMenu()
	// 生成菜单树
	//tree = GenMenuTree(nil, menus)
	return menus
}

// 生成菜单树
func GenMenuTree(parent *MenuTreeResp, menus []system.SysMenu) []MenuTreeResp {
	tree := make(MenuTreeRespList, 0)
	// 转为MenuTreeResponseStruct
	var resp []MenuTreeResp
	utils.Struct2StructByJson(menus, &resp)
	// parentId默认为0, 表示根菜单
	var parentId uint
	if parent != nil {
		parentId = parent.Id
	}
	for _, menu := range resp {
		// 父菜单编号一致
		if menu.ParentId == parentId {
			// 递归获取子菜单
			menu.Children = GenMenuTree(&menu, menus)
			// 加入菜单树
			tree = append(tree, menu)
		}
	}
	// 排序
	sort.Sort(tree)
	return tree
}

// 创建菜单
func CreateMenu(req *CreateMenuReq) (err error) {
	var menu system.SysMenu
	utils.Struct2StructByJson(req, &menu)
	// 创建数据
	err = common.Mysql.Create(&menu).Error
	return
}

// 更新菜单
func UpdateMenuById(id uint, req UpdateMenuReq) (err error) {
	if id == req.ParentId {
		return errors.New("不能自关联")
	}
	var oldMenu system.SysMenu
	query := common.Mysql.Table(oldMenu.TableName()).Where("id = ?", id).First(&oldMenu)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	// 比对增量字段,使用map确保gorm可更新零值
	var m map[string]interface{}
	utils.CompareDifferenceStructByJson(oldMenu, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除菜单
func DeleteMenuByIds(ids []uint) (err error) {
	var menu system.SysMenu
	// 先解除父级关联
	err = common.Mysql.Table(menu.TableName()).Where("parent_id IN (?)", ids).Update("parent_id", 0).Error
	if err != nil {
		return err
	}
	// 再删除
	err = common.Mysql.Where("id IN (?)", ids).Delete(&menu).Error
	if err != nil {
		return err
	}
	return
}

// 获取全部菜单, 非菜单树
func getAllMenu() []system.SysMenu {
	menus := make([]system.SysMenu, 0)
	// 查询所有菜单
	common.Mysql.Order("sort").Find(&menus)
	return menus
}
