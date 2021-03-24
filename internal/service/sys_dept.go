package service

import (
	"errors"
	"fmt"
	"go-xops/assets/system"
	"go-xops/internal/request"
	"go-xops/internal/response"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// 获取所有部门信息
func (s *MysqlService) GetDepts(req *request.DeptListReq) []system.SysDept {
	depts := make([]system.SysDept, 0)
	db := common.Mysql
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	status := req.Status
	if status != nil {
		if *status {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	db.Order("sort").Find(&depts)
	return depts
}

// 生成部门树
func GenDeptTree(parent *response.DeptTreeResp, depts []system.SysDept) []response.DeptTreeResp {
	tree := make(response.DeptTreeResppList, 0)
	var resp []response.DeptTreeResp
	utils.Struct2StructByJson(depts, &resp)
	// parentId默认为0, 表示根菜单
	var parentId uint
	if parent != nil {
		parentId = parent.Id
	}
	for _, dept := range resp {
		// 父菜单编号一致
		if dept.ParentId == parentId {
			// 递归获取子菜单
			dept.Children = GenDeptTree(&dept, depts)
			// 加入菜单树
			tree = append(tree, dept)
		}
	}
	// 排序
	sort.Sort(tree)
	return tree
}

// 创建部门
func (s *MysqlService) CreateDept(req *request.CreateDeptReq) (err error) {
	var dept system.SysDept
	utils.Struct2StructByJson(req, &dept)
	// 创建数据
	err = s.db.Create(&dept).Error
	return
}

// 更新部门
func (s *MysqlService) UpdateDeptById(id uint, req request.UpdateDeptReq) (err error) {
	if id == req.ParentId {
		return errors.New("不能自关联")
	}
	var oldDept system.SysDept
	query := s.db.Table(oldDept.TableName()).Where("id = ?", id).First(&oldDept)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	// 比对增量字段,使用map确保gorm可更新零值
	var m map[string]interface{}
	utils.CompareDifferenceStructByJson(oldDept, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除部门
func (s *MysqlService) DeleteDeptByIds(ids []uint) (err error) {
	var dept system.SysDept
	// 先解除父级关联
	err = s.db.Table(dept.TableName()).Where("parent_id IN (?)", ids).Update("parent_id", 0).Error
	if err != nil {
		return err
	}
	// 再删除
	err = s.db.Where("id IN (?)", ids).Delete(&dept).Error
	if err != nil {
		return err
	}
	return
}
