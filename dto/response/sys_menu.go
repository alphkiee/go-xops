package response

// 菜单树信息响应, 字段含义见models.SysMenu
type MenuTreeResp struct {
	Id       uint           `json:"id"`
	ParentId uint           `json:"parent_id"`
	Name     string         `json:"name"`
	Icon     string         `json:"icon"`
	Path     string         `json:"path"`
	Creator  string         `json:"creator"`
	Sort     int            `json:"sort"`
	Status   bool           `json:"status"`
	Children []MenuTreeResp `json:"children,omitempty"` //tag:omitempty 为空的值不显示
}
type MenuTreeRespList []MenuTreeResp

func (hs MenuTreeRespList) Len() int {
	return len(hs)
}
func (hs MenuTreeRespList) Less(i, j int) bool {
	return hs[i].Sort < hs[j].Sort // 按Sort从小到大排序
}

func (hs MenuTreeRespList) Swap(i, j int) {
	hs[i], hs[j] = hs[j], hs[i]
}

