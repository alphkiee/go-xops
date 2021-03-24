package cmdb

import (
	"go-xops/assets"
)

// Host 主机表
type Host struct {
	assets.Model
	HostName   string `gorm:"comment:'主机名';size:128" json:"host_name"`
	IP         string `gorm:"comment:'主机地址';size:128" json:"ip"`
	Port       string `gorm:"comment:'SSH端口';size:64" json:"port"`
	OsVersion  string `gorm:"comment:'系统版本';size:128" json:"os_version"`
	HostType   string `gorm:"comment:'主机类型';size:64" json:"host_type"`
	AuthType   string `gorm:"comment:'认证类型'" json:"auth_type"`
	User       string `gorm:"comment:'认证用户';size:64" json:"user"`
	Password   string `gorm:"comment:'认证密码';size:64" json:"password"`
	PrivateKey string `gorm:"comment:'秘钥';size:128" json:"privatekey"`
	Creator    string `gorm:"comment:'创建人';size:64" json:"creator"`
}

// TableName ...
func (m Host) TableName() string {
	return m.Model.TableName("cmdb_host")
}
