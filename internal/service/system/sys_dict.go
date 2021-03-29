package system

import (
	"errors"
	"fmt"
	"go-xops/assets/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"gorm.io/gorm"
)

// 创建字典结构体
type CreateDictReq struct {
	Key      string `json:"key" validate:"required"`
	Value    string `json:"value" validate:"required"`
	Desc     string `json:"desc"`
	ParentId uint   `json:"parent_id"`
	Creator  string `json:"creator"`
}

// 修改字典
type UpdateDictReq struct {
	Key      string `json:"key" validate:"required"`
	Value    string `json:"value" validate:"required"`
	Desc     string `json:"desc"`
	ParentId uint   `json:"parent_id"`
	Status   *bool  `json:"status"`
}

type DictListReq struct {
	Key     string `json:"key" form:"key"`
	Value   string `json:"value" form:"value"`
	Desc    string `json:"desc" form:"desc"`
	Creator string `json:"creator" form:"creator"`
	Status  *bool  `json:"status" form:"status"`
	TypeKey string `json:"type_key" form:"type_key"`
}

type DictTreeResp struct {
	Id       uint           `json:"id"`
	ParentId uint           `json:"parent_id"`
	Key      string         `json:"key"`
	Value    string         `json:"value"`
	Desc     string         `json:"desc"`
	Creator  string         `json:"creator"`
	Status   bool           `json:"status"`
	Children []DictTreeResp `json:"children,omitempty"` //tag:omitempty 为空的值不显示
}

// 获取所有字典信息
func GetDicts(req *DictListReq) []system.SysDict {
	Dicts := make([]system.SysDict, 0)
	db := common.Mysql
	typeKey := strings.TrimSpace(req.TypeKey)
	if typeKey != "" {
		var dist system.SysDict
		db = db.Preload("Dicts", "status = ?", true).Where("`key` LIKE ?", fmt.Sprintf("%%%s%%", typeKey)).First(&dist)
		return dist.Dicts
	}
	key := strings.TrimSpace(req.Key)
	if key != "" {
		db = db.Where("key LIKE ?", fmt.Sprintf("%%%s%%", key))
	}
	value := strings.TrimSpace(req.Value)
	if value != "" {
		db = db.Where("value LIKE ?", fmt.Sprintf("%%%s%%", value))
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
	db.Find(&Dicts)
	return Dicts
}

// 生成字典树
func GenDictTree(parent *DictTreeResp, Dicts []system.SysDict) []DictTreeResp {
	tree := make([]DictTreeResp, 0)
	var resp []DictTreeResp
	utils.Struct2StructByJson(Dicts, &resp)
	// parentId默认为0, 表示根菜单
	var parentId uint
	if parent != nil {
		parentId = parent.Id
	}
	for _, Dict := range resp {
		// 父菜单编号一致
		if Dict.ParentId == parentId {
			// 递归获取子菜单
			Dict.Children = GenDictTree(&Dict, Dicts)
			// 加入菜单树
			tree = append(tree, Dict)
		}
	}
	return tree
}

// 创建字典
func CreateDict(req *CreateDictReq) (err error) {
	var Dict system.SysDict
	utils.Struct2StructByJson(req, &Dict)
	// 创建数据
	err = common.Mysql.Create(&Dict).Error
	return
}

// 更新字典
func UpdateDictById(id uint, req UpdateDictReq) (err error) {
	if id == req.ParentId {
		return errors.New("不能自关联")
	}
	var oldDict system.SysDict
	query := common.Mysql.Table(oldDict.TableName()).Where("id = ?", id).First(&oldDict)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	// 比对增量字段,使用map确保gorm可更新零值
	var m map[string]interface{}
	utils.CompareDifferenceStructByJson(oldDict, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除字典
func DeleteDictByIds(ids []uint) (err error) {
	var Dict system.SysDict
	// 先解除父级关联
	err = common.Mysql.Table(Dict.TableName()).Where("parent_id IN (?)", ids).Update("parent_id", 0).Error
	if err != nil {
		return err
	}
	// 再删除
	err = common.Mysql.Where("id IN (?)", ids).Delete(&Dict).Error
	if err != nil {
		return err
	}
	return
}
