package cmd

import (
	"errors"
	"fmt"
	"go-xops/assets/cmdb"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HostListReq struct {
	Hostname      string `json:"hostname" form:"host_name"`
	IP            string `json:"ip" form:"ip" validate:"required"`
	Port          string `json:"port" form:"port"`
	User          string `json:"user" form:"user"`
	OsVersion     string `json:"os_version" form:"os_version"`
	HostType      string `json:"host_type" form:"host_type"`
	AuthType      string `json:"auth_type" form:"auth_type" validate:"required"`
	Password      string `json:"password" form:"password"`
	PrivateKey    string `json:"privatekey" form:"privatekey"`
	KeyPassphrase string `json:"key_passphrase"`
	Creator       string `json:"creator" form:"creator"`
	common.PageInfo
}

// CreateHostReq 创建主机
type CreateHostReq struct {
	HostName   string `json:"host_name" form:"host_name"`
	IP         string `json:"ip" form:"ip_address" validate:"required"`
	HostType   string `json:"host_type" form:"host_type"`
	Port       string `json:"port" form:"port"`
	AuthType   string `json:"auth_type" form:"auth_type" validate:"required"`
	User       string `json:"user" form:"user"`
	Password   string `json:"password" form:"password"`
	OsVersion  string `json:"os_version" form:"os_version"`
	PrivateKey string `json:"privatekey" form:"privatekey"`
	Creator    string `json:"creator" form:"creator"`
}

// HostRep ...
type HostRep struct {
	HostName  string `json:"host_name"`
	IP        string `json:"ip"`
	OsVersion string `json:"os_version"`
	AuthType  string `json:"auth_type"`
	Creator   string `json:"creator"`
}

// GetHosts 获取主机
func GetHosts(req *HostListReq) ([]cmdb.Host, error) {
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
func CreateHost(req *CreateHostReq) (err error) {
	var host cmdb.Host
	utils.Struct2StructByJson(req, &host)

	err = common.Mysql.Create(&host).Error
	return
}

// 更新主机
func UpdateHostById(id uint, req gin.H) (err error) {
	var oldHost cmdb.Host
	query := common.Mysql.Where("id = ?", id).First(&oldHost)
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
func DeleteHostsById(id []uint) (err error) {
	var host cmdb.Host
	err = common.Mysql.Where("id in ?", id).Delete(&host).Error
	return err
}

// 根据ID来获取主机
func GetHostByid(id uint) (cmdb.Host, error) {
	var host cmdb.Host

	err := common.Mysql.Where("id = ?", id).Find(&host).Error
	return host, err
}

// GetHostByIds 获取主机根据传入过来的ID
func GetHostByIds(ids []uint) ([]cmdb.Host, error) {
	var host []cmdb.Host
	err := common.Mysql.Where("id in ?", ids).Find(&host).Error
	return host, err
}
