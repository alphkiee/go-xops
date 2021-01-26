package service

import (
	"fmt"
	"go-xops/dto/request"
	"go-xops/models/system"
	"go-xops/pkg/common"
	"strings"
)

// 获取操作日志
func (s *MysqlService) GetOperLogs(req *request.OperLogListReq) ([]system.SysOperLog, error) {
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
func (s *MysqlService) DeleteOperationLogByIds(ids []uint) (err error) {
	return s.db.Where("id IN (?)", ids).Delete(system.SysOperLog{}).Error
}
