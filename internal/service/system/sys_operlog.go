package system

import (
	"fmt"
	"go-xops/assets"
	"go-xops/assets/system"
	"go-xops/pkg/common"
	"strings"
	"time"
)

// 获取操作日志列表结构体
type OperLogListReq struct {
	Name            string `json:"name" form:"name"`
	Method          string `json:"method" form:"method"`
	Path            string `json:"path" form:"path"`
	Username        string `json:"username" form:"username"`
	Ip              string `json:"ip" form:"ip"`
	common.PageInfo        // 分页参数
}

// 翻译需要校验的字段名称
func (s OperLogListReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Status"] = "响应状态码"
	return m
}

// 创建操作日志结构体
type CreateOperLogRequestStruct struct {
	Path       string        `json:"path"`
	Method     string        `json:"method"`
	Params     string        `json:"params"`
	Body       string        `json:"body"`
	Data       string        `json:"data"`
	Status     int           `json:"status"`
	Username   string        `json:"username"`
	Ip         string        `json:"ip"`
	IpLocation string        `json:"ip_location"`
	Latency    time.Duration `json:"latency"`
	UserAgent  string        `json:"user_agent"`
}

type OperationLogListResp struct {
	Id         uint             `json:"id"`
	Name       string           `json:"name"`
	Path       string           `json:"path"`
	Method     string           `json:"method"`
	Params     string           `json:"params"`
	Body       string           `json:"body"`
	Data       string           `json:"data"`
	Status     int              `json:"status"`
	Username   string           `json:"username"`
	Ip         string           `json:"ip"`
	IpLocation string           `json:"ip_location"`
	Latency    time.Duration    `json:"latency"`
	UserAgent  string           `json:"user_agent"`
	CreatedAt  assets.LocalTime `json:"created_at"`
}

// 获取操作日志
func GetOperLogs(req *OperLogListReq) ([]system.SysOperLog, error) {
	var err error
	list := make([]system.SysOperLog, 0)
	query := common.Mysql
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
	username := strings.TrimSpace(req.Username)
	if path != "" {
		query = query.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		query = query.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	query = query.Order("id DESC")
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

// 批量删除操作日志
func DeleteOperationLogByIds(ids []uint) (err error) {
	return common.Mysql.Where("id IN (?)", ids).Delete(system.SysOperLog{}).Error
}
