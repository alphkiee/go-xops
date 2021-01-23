package request

import (
	"go-xops/dto/response"
)

// User login structure
type RegisterAndLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 修改密码结构体
type ChangePwdReq struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required"`
}

// 用户列表请求结构体
type UserListReq struct {
	Username          string `json:"username" form:"username"`
	Mobile            string `json:"mobile" form:"mobile"`
	Name              string `json:"name" form:"name"`
	Status            *bool  `json:"status" form:"status"`
	Creator           string `json:"creator" form:"creator"`
	DeptId            uint   `json:"dept_id" form:"dept_id"`
	response.PageInfo        // 分页参数
}

// 创建用户结构体
type CreateUserReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"`
	DeptId   uint   `json:"dept_id"`
	RoleId   uint   `json:"role_id" validate:"required"`
	Creator  string `json:"creator"`
}

// 修改用户基本信息结构体
type UpdateUserBaseInfoReq struct {
	Mobile string `json:"mobile"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email"`
}

// 修改用户结构体
type UpdateUserReq struct {
	Mobile   string `json:"mobile"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   uint   `json:"role_id"`
	DeptId   uint   `json:"dept_id"`
	Status   *bool  `json:"status"`
}

// 翻译需要校验的字段名称
func (s CreateUserReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Username"] = "用户名"
	m["Password"] = "用户密码"
	m["Name"] = "姓名"
	m["RoleId"] = "角色ID"
	return m
}
func (s ChangePwdReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["OldPassword"] = "旧密码"
	m["NewPassword"] = "新密码"
	return m
}

func (s UpdateUserReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "姓名"
	return m
}

func (s UpdateUserBaseInfoReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "姓名"
	return m
}
