package system

import (
	"errors"
	"fmt"
	"go-xops/assets"
	"go-xops/assets/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LoginResp struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	Expires  string `json:"expires"`
}

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

type UpdateUserBaseInfoReq struct {
	Mobile string `json:"mobile"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email"`
}

type UpdateUserReq struct {
	Mobile   string `json:"mobile"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   uint   `json:"role_id"`
	DeptId   uint   `json:"dept_id"`
	Status   *bool  `json:"status"`
}

type UserListReq struct {
	Username string `json:"username" form:"username"`
	Mobile   string `json:"mobile" form:"mobile"`
	Name     string `json:"name" form:"name"`
	Status   *bool  `json:"status" form:"status"`
	Creator  string `json:"creator" form:"creator"`
	DeptId   uint   `json:"dept_id" form:"dept_id"`
	common.PageInfo
}

type ChangePwdReq struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required"`
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

// 用户列表信息响应, 字段含义见assets.SysUser
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
	CreatedAt assets.LocalTime `json:"created_at"`
	UpdatedAt assets.LocalTime `json:"updated_at"`
}

func LoginCheck(username, password string) (*LoginResp, error) {
	var u system.SysUser
	err := common.Mysql.Preload("Role", "status = ?", true).Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, errors.New(common.LoginCheckErrorMsg)
	}
	if !*u.Status {
		return nil, errors.New(common.UserForbiddenMsg)
	}
	// 校验密码
	if ok := utils.ComparePwd(password, u.Password); !ok {
		return nil, errors.New(common.LoginCheckErrorMsg)
	}
	var loginInfo LoginResp
	utils.Struct2StructByJson(u, &loginInfo)
	return &loginInfo, err
}

func GetUserById(id uint) (system.SysUser, error) {
	var user system.SysUser
	var err error
	err = common.Mysql.Preload("Role", "status = ?", true).Preload("Dept", "status = ?", true).Where("id = ?", id).First(&user).Error
	return user, err
}

// 检查用户是否已存在
func CheckUser(username string) error {
	var user system.SysUser
	var err error
	common.Mysql.Where("username = ?", username).First(&user)
	if user.Id != 0 {
		err = errors.New("用户名已存在")
	}
	return err
}

// 创建用户
func CreateUser(req *CreateUserReq) (err error) {
	var user system.SysUser
	err = CheckUser(req.Username)
	if err != nil {
		return
	}
	utils.Struct2StructByJson(req, &user)
	user.Password = utils.GenPwd(req.Password)
	err = common.Mysql.Create(&user).Error
	return
}

// UpdateUserBaseInfoById ...更新用户基本信息
func UpdateUserBaseInfoById(id uint, req UpdateUserBaseInfoReq) (err error) {
	var oldUser system.SysUser
	query := common.Mysql.Table(oldUser.TableName()).Where("id = ?", id).First(&oldUser)
	if query.Error == gorm.ErrInvalidField {
		return errors.New("记录不存在")
	}
	var m system.SysUser
	utils.CompareDifferenceStructByJson(oldUser, req, &m)
	// 更新指定列
	err = query.Updates(m).Error
	return
}

// UpdateUserById ...更新用户
func UpdateUserById(id uint, req UpdateUserReq) (err error) {
	var oldUser system.SysUser
	query := common.Mysql.Table(oldUser.TableName()).Where("id = ?", id).First(&oldUser)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	logrus.Printf("req= %v", req)
	password := ""
	if strings.TrimSpace(req.Password) != "" {
		password = utils.GenPwd(req.Password)
	}
	var m map[string]interface{}
	utils.CompareDifferenceStructByJson(oldUser, req, &m)
	delete(m, "password")
	if password != "" {
		err = query.Update("password", password).Updates(m).Error
	} else {
		err = query.Updates(m).Error
	}
	return
}

// 获取用户
func GetUsers(req *UserListReq) ([]system.SysUser, error) {
	var err error
	list := make([]system.SysUser, 0)
	db := common.Mysql
	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	status := req.Status
	if status != nil {
		if *status {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.All {
			err = db.Preload("Role", "status = ?", true).Preload("Dept", "status = ?", true).Find(&list).Error
		} else {
			limit, offset := req.GetLimit()
			err = db.Preload("Role", "status = ?", true).Preload("Dept", "status = ?", true).Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 批量删除用户
func DeleteUserByIds(ids []uint) (err error) {
	var user system.SysUser
	return common.Mysql.Where("id IN (?)", ids).Delete(&user).Error
}
