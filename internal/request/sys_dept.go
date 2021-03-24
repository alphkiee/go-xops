package request

// 创建部门结构体
type CreateDeptReq struct {
	Name      string `json:"name" validate:"required"`
	Sort       int    `json:"sort"`
	ParentId   uint   `json:"parent_id"`
	Creator    string `json:"creator"`
}

// 修改部门
type UpdateDeptReq struct {
	Name     string `json:"name" validate:"required"`
	Status   *bool  `json:"status"`
	Sort       int    `json:"sort"`
	ParentId   uint   `json:"parent_id"`
}

type DeptListReq struct {
	Name              string `json:"name" form:"name"`
	Creator           string `json:"creator" form:"creator"`
	Status            *bool  `json:"status" form:"status"`
}

// 翻译需要校验的字段名称
func (s CreateDeptReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "部门名称"
	return m
}

// 翻译需要校验的字段名称
func (s UpdateDeptReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "部门名称"
	return m
}
