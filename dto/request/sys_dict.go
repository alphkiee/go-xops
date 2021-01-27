package request

// 创建字典结构体
type CreateDictReq struct {
	Key      string `json:"key" validate:"required"`
	Value    string `json:"value" validate:"required"`
	Desc     string `json:"desc"`
	ParentId uint   `json:"parent_id"`
	Creator  string `json:"creator"`
}

// 修改字典
type UpdateDictReq struct {
	Key      string `json:"key" validate:"required"`
	Value    string `json:"value" validate:"required"`
	Desc     string `json:"desc"`
	ParentId uint   `json:"parent_id"`
	Status   *bool  `json:"status"`
}

type DictListReq struct {
	Key     string `json:"key" form:"key"`
	Value   string `json:"value" form:"value"`
	Desc    string `json:"desc" form:"desc"`
	Creator string `json:"creator" form:"creator"`
	Status  *bool  `json:"status" form:"status"`
	TypeKey string `json:"type_key" form:"type_key"`
}

// 翻译需要校验的字段名称
func (s CreateDictReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Key"] = "字典Key"
	m["Value"] = "字典Value"
	return m
}

// 翻译需要校验的字段名称
func (s UpdateDictReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Key"] = "字典Key"
	m["Value"] = "字典Value"
	return m
}
