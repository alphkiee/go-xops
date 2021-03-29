package system

import (
	"errors"
	"fmt"
	"go-xops/assets"
	"go-xops/assets/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取接口列表结构体
type ApiListReq struct {
	Name            string `json:"name" form:"name"`
	Method          string `json:"method" form:"method"`
	Path            string `json:"path" form:"path"`
	Category        string `json:"category" form:"category"`
	Creator         string `json:"creator" form:"creator"`
	Tree            bool   `json:"tree" form:"tree"`
	common.PageInfo        // 分页参数
}

// 创建接口结构体
type CreateApiReq struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Method   string `json:"method" validate:"required"`
	Path     string `json:"path" validate:"required"`
	Category string `json:"category" validate:"required"`
	Creator  string `json:"creator" form:"creator"`
	Desc     string `json:"desc"`
}

// ApiListResp ...接口信息响应, 字段含义见models
type ApiListResp struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	Method    string           `json:"method"`
	Path      string           `json:"path"`
	Category  string           `json:"category"`
	Creator   string           `json:"creator"`
	Desc      string           `json:"desc"`
	CreatedAt assets.LocalTime `json:"created_at"`
}

type ApiTreeResp struct {
	Category string        `json:"category"` // 分组名称
	Children []ApiListResp `json:"children"` // 前端以树形图结构展示, 这里用children表示
}

func GetApis(req *ApiListReq) ([]system.SysApi, error) {
	var err error
	list := make([]system.SysApi, 0)
	query := common.Mysql.Table(new(system.SysApi).TableName())
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	method := strings.TrimSpace(req.Method)
	if method != "" {
		query = query.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		query = query.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		query = query.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}

	// 查询条数
	err = query.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.All {
			// 不使用分页
			err = query.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = query.Limit(limit).Offset(offset).Find(&list).Error
		}
	}

	return list, err
}

// 创建接口
func CreateApi(req *CreateApiReq) (err error) {
	var api system.SysApi
	utils.Struct2StructByJson(req, &api)
	// 创建数据
	err = common.Mysql.Create(&api).Error
	return
}

// 更新接口
func UpdateApiById(id uint, req gin.H) (err error) {
	var oldApi system.SysApi
	query := common.Mysql.Table(oldApi.TableName()).Where("id = ?", id).First(&oldApi)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	var m system.SysApi
	utils.CompareDifferenceStructByJson(oldApi, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除接口
func DeleteApiByIds(ids []uint) (err error) {

	return common.Mysql.Where("id IN (?)", ids).Delete(system.SysApi{}).Error
}
