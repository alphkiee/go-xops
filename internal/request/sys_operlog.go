package request

import (
	"go-xops/internal/response"
	"time"
)

// 获取操作日志列表结构体
type OperLogListReq struct {
	Name              string `json:"name" form:"name"`
	Method            string `json:"method" form:"method"`
	Path              string `json:"path" form:"path"`
	Username          string `json:"username" form:"username"`
	Ip                string `json:"ip" form:"ip"`
	response.PageInfo        // 分页参数
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
