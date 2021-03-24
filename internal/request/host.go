package request

import "go-xops/internal/response"

// HostListReq 主机列表
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
	response.PageInfo
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
	//KeyPassphrase string `json:"key_passphrase"`
	Creator string `json:"creator" form:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateHostReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["IpAddress"] = "主机地址"
	m["AuthType"] = "认证类型"
	return m
}
