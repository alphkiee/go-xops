package service

import (
	"errors"
	"fmt"
	"go-xops/assets/cmdb"
	"go-xops/internal/request"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetHosts 获取主机
func (s *MysqlService) GetHosts(req *request.HostListReq) ([]cmdb.Host, error) {
	hosts := make([]cmdb.Host, 0)
	db := common.Mysql

	host_name := strings.TrimSpace(req.Hostname)
	if host_name != "" {
		db = db.Where("host_name LIKE ?", fmt.Sprintf("%%%%s%%", host_name))
	}

	ip := strings.TrimSpace(req.IP)
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}

	host_type := strings.TrimSpace(req.HostType)
	if host_type != "" {
		db.Where("host_type LIKE ?", fmt.Sprintf("%%%s%%", host_type))
	}

	auth_type := strings.TrimSpace(req.AuthType)
	if auth_type != "" {
		db.Where("auth_type LIKE ?", fmt.Sprintf("%%%s%%", auth_type))
	}
	err := db.Find(&hosts).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.All {
			err = db.Find(&hosts).Error
		} else {
			limit, offset := req.GetLimit()
			err = db.Limit(limit).Offset(offset).Find(&hosts).Error
		}
	}
	return hosts, err
}

// 添加主机
func (s *MysqlService) CreateHost(req *request.CreateHostReq) (err error) {
	var host cmdb.Host
	utils.Struct2StructByJson(req, &host)

	err = s.db.Create(&host).Error
	return
}

// 更新主机
func (s *MysqlService) UpdateHostById(id uint, req gin.H) (err error) {
	var oldHost cmdb.Host
	query := s.db.Where("id = ?", id).First(&oldHost)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}

	// 比对增量字段
	var m cmdb.Host
	utils.CompareDifferenceStructByJson(oldHost, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// 批量删除主机
func (s *MysqlService) DeleteHostsById(id []uint) (err error) {
	var host cmdb.Host
	err = s.db.Where("id in ?", id).Delete(&host).Error
	return err
}

// 根据ID来获取主机
func (s *MysqlService) GetHostByid(id uint) (cmdb.Host, error) {
	var host cmdb.Host

	err := s.db.Where("id = ?", id).Find(&host).Error
	return host, err
}

// GetHostByIds 获取主机根据传入过来的ID
func (s *MysqlService) GetHostByIds(ids []uint) ([]cmdb.Host, error) {
	var host []cmdb.Host
	err := s.db.Where("id in ?", ids).Find(&host).Error
	return host, err
}
