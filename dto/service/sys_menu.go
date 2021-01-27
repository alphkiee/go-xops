package service

import (
	"errors"
	"go-xops/dto/request"
	"go-xops/dto/response"
	"go-xops/models/system"
	"go-xops/pkg/utils"
	"sort"

	"gorm.io/gorm"
)

// 获取用户菜单的切片
func (s *MysqlService) GetUserMenuList(roleId uint) ([]system.SysMenu, error) {
	//tree := make([]models.SysMenu, 0)
	var role system.SysRole
	err := s.db.Table(new(system.SysRole).TableName()).Preload("Menus", "status = ?", true).Where("id = ?", roleId).Find(&role).Error
	menus := make([]system.SysMenu, 0)
	if err != nil {
		return menus, err
	}
	if role.Keyword == "admin" {
		err = s.db.Find(&menus).Error
		return menus, err
	}
	// 生成菜单树
	//tree = GenMenuTree(nil, role.Menus)
	return role.Menus, nil
}

// 获取所有菜单
func (s *MysqlService) GetMenus() []system.SysMenu {
	//tree := make([]models.SysMenu, 0)
	menus := s.getAllMenu()
	// 生成菜单树
	//tree = GenMenuTree(nil, menus)
	return menus
}

// 生成菜单树
func GenMenuTree(parent *response.MenuTreeResp, menus []system.SysMenu) []response.MenuTreeResp {
	tree := make(response.MenuTreeRespList, 0)
	// 转为MenuTreeResponseStruct
	var resp []response.MenuTreeResp
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
func (s *MysqlService) CreateMenu(req *request.CreateMenuReq) (err error) {
	var menu system.SysMenu
	utils.Struct2StructByJson(req, &menu)
	// 创建数据
	err = s.db.Create(&menu).Error
	return
}

// 更新菜单
func (s *MysqlService) UpdateMenuById(id uint, req request.UpdateMenuReq) (err error) {
	if id == req.ParentId {
		return errors.New("不能自关联")
	}
	var oldMenu system.SysMenu
	query := s.db.Table(oldMenu.TableName()).Where("id = ?", id).First(&oldMenu)
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
func (s *MysqlService) DeleteMenuByIds(ids []uint) (err error) {
	var menu system.SysMenu
	// 先解除父级关联
	err = s.db.Table(menu.TableName()).Where("parent_id IN (?)", ids).Update("parent_id", 0).Error
	if err != nil {
		return err
	}
	// 再删除
	err = s.db.Where("id IN (?)", ids).Delete(&menu).Error
	if err != nil {
		return err
	}
	return
}

// 获取全部菜单, 非菜单树
func (s *MysqlService) getAllMenu() []system.SysMenu {
	menus := make([]system.SysMenu, 0)
	// 查询所有菜单
	s.db.Order("sort").Find(&menus)
	return menus
}
