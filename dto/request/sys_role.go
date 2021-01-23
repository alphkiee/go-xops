package request

import "go-xops/dto/response"

// 获取角色列表结构体
type RoleListReq struct {
	Name              string `json:"name" form:"name"`
	Keyword           string `json:"keyword" form:"keyword"`
	Status            *bool  `json:"status" form:"status"`
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// 更新角色权限的结构体
type UpdateRolePermsReq struct {
	MenusId []uint `json:"menus_id" form:"menus_id"` // 传多个id
	ApisId  []uint `json:"apis_id" form:"apis_id"`   // 传多个id
}

// 创建角色结构体
type CreateRoleReq struct {
	Name    string `json:"name" validate:"required"`
	Keyword string `json:"keyword" validate:"required"`
	Desc    string `json:"desc"`
	Creator string `json:"creator"`
}

// 翻译需要校验的字段名称
func (s CreateRoleReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "角色名称"
	m["Keyword"] = "角色关键字"
	return m
}
