package service

import (
	"errors"
	"fmt"
	"go-xops/assets/system"
	"go-xops/internal/request"
	"go-xops/internal/response"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 登录校验,返回指定信息
func (s *MysqlService) LoginCheck(username string, password string) (*response.LoginResp, error) {
	var u system.SysUser
	// 查询用户及其角色
	err := s.db.Preload("Role", "status = ?", true).Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	if !*u.Status {
		return nil, errors.New(response.UserForbiddenMsg)
	}
	// 校验密码
	if ok := utils.ComparePwd(password, u.Password); !ok {
		return nil, errors.New(response.LoginCheckErrorMsg)
	}
	var loginInfo response.LoginResp
	utils.Struct2StructByJson(u, &loginInfo)
	loginInfo.CurrentAuthority = []string{u.Role.Keyword}
	return &loginInfo, err
}

// 获取单个用户
func (s *MysqlService) GetUserById(id uint) (system.SysUser, error) {
	var user system.SysUser
	var err error
	err = s.db.Preload("Role", "status = ?", true).Preload("Dept", "status = ?", true).Where("id = ?", id).First(&user).Error
	return user, err
}

// 检查用户是否已存在
func (s *MysqlService) CheckUser(username string) error {
	var user system.SysUser
	var err error
	s.db.Where("username = ?", username).First(&user)
	if user.Id != 0 {
		err = errors.New("用户名已存在")
	}
	return err
}

// 创建用户
func (s *MysqlService) CreateUser(req *request.CreateUserReq) (err error) {
	var user system.SysUser
	err = s.CheckUser(req.Username)
	if err != nil {
		return
	}
	utils.Struct2StructByJson(req, &user)
	// 将初始密码转为密文
	user.Password = utils.GenPwd(req.Password)
	// 创建数据
	err = s.db.Create(&user).Error
	return
}

// UpdateUserBaseInfoById ...更新用户基本信息
func (s *MysqlService) UpdateUserBaseInfoById(id uint, req request.UpdateUserBaseInfoReq) (err error) {
	var oldUser system.SysUser
	query := s.db.Table(oldUser.TableName()).Where("id = ?", id).First(&oldUser)
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
func (s *MysqlService) UpdateUserById(id uint, req request.UpdateUserReq) (err error) {
	var oldUser system.SysUser
	query := s.db.Table(oldUser.TableName()).Where("id = ?", id).First(&oldUser)
	// query := s.db.First(&oldUser).Where("id = ?", id)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}
	logrus.Printf("req= %v", req)
	password := ""
	// 填写了新密码
	if strings.TrimSpace(req.Password) != "" {
		password = utils.GenPwd(req.Password)
	}
	// 比对增量字段,使用map确保gorm可更新零值
	var m map[string]interface{}
	utils.CompareDifferenceStructByJson(oldUser, req, &m)
	delete(m, "password")
	if password != "" {
		// 更新密码
		err = query.Update("password", password).Updates(m).Error
	} else {
		// 更新指定列
		err = query.Updates(m).Error
	}
	return
}

// 获取用户
func (s *MysqlService) GetUsers(req *request.UserListReq) ([]system.SysUser, error) {
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
	// 查询条数
	err = db.Find(&list).Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.All {
			// 不使用分页
			err = db.Preload("Role", "status = ?", true).Preload("Dept", "status = ?", true).Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = db.Preload("Role", "status = ?", true).Preload("Dept", "status = ?", true).Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 批量删除用户
func (s *MysqlService) DeleteUserByIds(ids []uint) (err error) {
	var user system.SysUser
	return s.db.Where("id IN (?)", ids).Delete(&user).Error
}
