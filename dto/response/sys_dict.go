package response

// 字典树信息响应,

type DictTreeResp struct {
	Id       uint           `json:"id"`
	ParentId uint           `json:"parent_id"`
	Key     string         `json:"key"`
	Value     string         `json:"value"`
	Desc     string         `json:"desc"`
	Creator  string         `json:"creator"`
	Status   bool           `json:"status"`
	Children []DictTreeResp `json:"children,omitempty"` //tag:omitempty 为空的值不显示
}