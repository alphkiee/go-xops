package system

import (
	"go-xops/models"
	"time"
)

// 系统操作日志
type SysOperLog struct {
	models.Model
	Name       string        `gorm:"comment:'接口名称';size:128" json:"name"`
	Path       string        `json:"path" gorm:"comment:'访问路径';size:128"`
	Method     string        `json:"method" gorm:"comment:'请求方式';size:128"`
	Body       string        `json:"body" gorm:"type:blob;comment:'请求主体(通过二进制存储节省空间)';size:128"`
	Data       string        `json:"data" gorm:"type:blob;comment:'响应数据(通过二进制存储节省空间)';size:128"`
	Status     int           `json:"status" gorm:"comment:'响应状态码'"`
	Username   string        `json:"username" gorm:"comment:'用户登录名';size:128"`
	Ip         string        `json:"ip" gorm:"comment:'Ip地址';size:128"`
	IpLocation string        `json:"ip_location" gorm:"comment:'Ip所在地';size:128"`
	Latency    time.Duration `json:"latency" gorm:"comment:'请求耗时(ms)'"`
	UserAgent  string        `json:"user_agent" gorm:"comment:'浏览器标识';size:128"`
}

func (m SysOperLog) TableName() string {
	return m.Model.TableName("sys_operlog")
}
