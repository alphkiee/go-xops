package response

import (
	"go-xops/models"
	"time"
)

// 接口信息响应, 字段含义见models
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
	CreatedAt  models.LocalTime `json:"created_at"`
}
