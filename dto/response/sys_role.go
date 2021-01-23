package response

import (
	"go-xops/models"
)

// 角色返回权限信息
type RolePermsResp struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
	MenusId []uint `json:"menus_id"`
	ApisId  []uint `json:"apis_id"`
}

// 角色信息响应, 字段含义见models
type RoleListResp struct {
	Id        uint             `json:"id"`
	Name      string           `json:"name"`
	Keyword   string           `json:"keyword"`
	Desc      string           `json:"desc"`
	Status    *bool            `json:"status"`
	Creator   string           `json:"creator"`
	CreatedAt models.LocalTime `json:"created_at"`
}
