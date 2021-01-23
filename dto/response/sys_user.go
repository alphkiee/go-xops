package response

import (
	"go-xops/models"
)

// User login response structure
type LoginResp struct {
	Id               uint     `json:"id"`
	Username         string   `json:"username"`
	Avatar           string   `json:"avatar"`
	Name             string   `json:"name"`
	Token            string   `json:"token"`            // jwt令牌
	Expires          string   `json:"expires"`          // 过期时间, 秒
	CurrentAuthority []string `json:"currentAuthority"` // 返回前端的权限数据
}

// 用户返回角色信息
type UserRoleResp struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Keyword string `json:"keyword"`
	Status  *bool  `json:"status"`
}

type UserDeptResp struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

// 用户信息响应
type UserInfoResp struct {
	Id       uint         `json:"id"`
	Username string       `json:"username"`
	Mobile   string       `json:"mobile"`
	Avatar   string       `json:"avatar"`
	Name     string       `json:"name"`
	Email    string       `json:"email"`
	Dept     UserDeptResp `json:"dept"`
	Role     UserRoleResp `json:"role"`
}

// 用户列表信息响应, 字段含义见models.SysUser
type UserListResp struct {
	Id        uint             `json:"id"`
	Username  string           `json:"username"`
	Mobile    string           `json:"mobile"`
	Avatar    string           `json:"avatar"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Status    *bool            `json:"status"`
	DeptId    uint             `json:"dept_id"`
	RoleId    uint             `json:"role_id"`
	Role      UserRoleResp     `json:"role"`
	Dept      UserDeptResp     `json:"dept"`
	Creator   string           `json:"creator"`
	CreatedAt models.LocalTime `json:"created_at"`
	UpdatedAt models.LocalTime `json:"updated_at"`
}
