package response

import (
	"go-xops/models"
)

// ApiListResp ...接口信息响应, 字段含义见models
type ApiListResp struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	Method    string           `json:"method"`
	Path      string           `json:"path"`
	Category  string           `json:"category"`
	Creator   string           `json:"creator"`
	Desc      string           `json:"desc"`
	CreatedAt models.LocalTime `json:"created_at"`
}

type ApiTreeResp struct {
	Category string        `json:"category"` // 分组名称
	Children []ApiListResp `json:"children"` // 前端以树形图结构展示, 这里用children表示
}
