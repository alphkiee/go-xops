package system

type SysCasbin struct {
	PType string `gorm:"size:100;comment:'策略类型'"`
	V0    string `gorm:"size:100;comment:'角色关键字'"`
	V1    string `gorm:"size:100;comment:'资源名称'"`
	V2    string `gorm:"size:100;comment:'请求类型'"`
	V3    string `gorm:"size:100"`
	V4    string `gorm:"size:100"`
	V5    string `gorm:"size:100"`
}

func (m SysCasbin) TableName() string {

	return "casbin_rule"
}

// SysRoleCasbin ...
type SysRoleCasbin struct {
	Keyword string `json:"keyword"`
	Method  string `json:"method"`
	Path    string `json:"path"`
}
