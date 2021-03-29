package system

import (
	"go-xops/assets/system"
	s "go-xops/internal/service/system"
	"go-xops/pkg/cache"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetCurrentUserFromCache ...
func GetCurrentUserFromCache(c *gin.Context) interface{} {
	user, exists := c.Get("user")
	var newUser system.SysUser
	if !exists {
		return newUser
	}
	u, _ := user.(s.LoginResp)
	// 创建缓存对象
	cache, err := cache.New(time.Second * 15)
	if err != nil {
		logrus.Error(err)
	}
	key := "user:" + u.Username
	cache.DBGetter = func() interface{} {
		newUser, _ = s.GetUserById(u.Id)
		return newUser
	}

	cache.GetCache(key)
	return newUser
}

// GetUserInfo doc
// @Summary Get /api/v1/user/info
// @Description 获取当前用户信息
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Router /api/v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 转为UserInfoResponseStruct, 隐藏部分字段
	var resp s.UserInfoResp
	utils.Struct2StructByJson(user, &resp)
	common.SuccessWithData(resp)
}

// CreateUser doc
// @Summary Post /api/v1/user/create
// @Description 创建用户
// @Produce json
// @Param data body request.CreateUserReq true "username, password, name, role_id"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/user/create [post]
func CreateUser(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	var req s.CreateUserReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	m := make(map[string]string, 0)
	m["Username"] = "用户名"
	m["Password"] = "密码"
	m["Name"] = "姓名"
	m["RoleId"] = "角色id"
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	req.Creator = user.(system.SysUser).Name
	err = s.CreateUser(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// GetUsers doc
// @Summary Get /api/v1/user/list
// @Description 列出所有用户
// @Produce json
// @Param username query string false "username"
// @Param mobile query string false "mobile"
// @Param name query string false "name"
// @Param status query string false "status"
// @Param creator query string false "creator"
// @Param dept_id query string false "dept_id"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/user/list [get]
func GetUsers(c *gin.Context) {
	// 绑定参数
	var req s.UserListReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	users, err := s.GetUsers(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	var respStruct []s.UserListResp
	utils.Struct2StructByJson(users, &respStruct)
	// 返回分页数据
	var resp common.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	common.SuccessWithData(resp)

}

// UpdateUserBaseInfoById doc
// @Summary Patch /api/v1/user/info/update/:userId
// @Description 根据user ID来更新用户基本信息
// @Produce json
// @Param userId path int true "userId"
// @Param data body request.UpdateUserBaseInfoReq true "mobile, name, email"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/user/info/update/{userId} [patch]
func UpdateUserBaseInfoById(c *gin.Context) {
	// 绑定参数
	var req s.UpdateUserBaseInfoReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["Mobile"] = "手机"
	m["Name"] = "姓名"
	m["Email"] = "邮箱"
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		common.FailWithMsg("用户编号不正确")
		return
	}
	err = s.UpdateUserBaseInfoById(userId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateUserById doc
// @Summary Patch /api/v1/user/update/:userId
// @Description 更新用户根据 user ID
// @Produce json
// @Param userId path int true "userId"
// @Param data body request.UpdateUserReq true "mobile, name, email, password"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/user/update/{userId} [patch]
func UpdateUserById(c *gin.Context) {
	// 绑定参数
	var req s.UpdateUserReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["Name"] = "姓名"
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		common.FailWithMsg("用户编号不正确")
		return
	}
	err = s.UpdateUserById(userId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// ChangePwd doc
// @Summary Put /api/v1/user/changePwd
// @Description 更改用户的密码
// @Produce json
// @Param data body request.ChangePwdReq true "old_password, new_password"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/user/changePwd [put]
func ChangePwd(c *gin.Context) {
	var msg string
	var req s.ChangePwdReq
	_ = c.ShouldBindJSON(&req)
	m := make(map[string]string, 0)
	m["OldPassword"] = "旧密码"
	m["NewPassword"] = "新密码"
	err := common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	user := GetCurrentUserFromCache(c)
	query := common.Mysql.Where("username = ?", user.(system.SysUser).Username).First(&user)
	err = query.Error
	if err != nil {
		msg = err.Error()
	} else {
		if ok := utils.ComparePwd(req.OldPassword, user.(system.SysUser).Password); !ok {
			msg = "原密码错误"
		} else {
			err = query.Update("password", utils.GenPwd(req.NewPassword)).Error
			if err != nil {
				msg = err.Error()
			}
		}
	}
	if msg != "" {
		common.FailWithMsg(msg)
		return
	}
	common.Success()
}

// DeleteUserByIds doc
// @Summary Delete /api/v1/user/delete
// @Description 根据ID批量删除用户
// @Produce json
// @Param data body request.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/user/delete [delete]
func DeleteUserByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = s.DeleteUserByIds(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UserAvatarUpload doc
// @Summary Post /api/v1/user/info/uploadImg
// @Description 上传头像
// @Produce json
// @Param avatar formData file true "avatar"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/user/info/uploadImg [post]
func UserAvatarUpload(c *gin.Context) {
	// 限制头像2MB(二进制移位xxxMB)
	err := c.Request.ParseMultipartForm(2 << 20)
	if err != nil {
		common.FailWithMsg("文件为空或图片大小超出最大值2MB")
		return
	}
	// 读取文件
	file, err := c.FormFile("avatar")
	if err != nil {
		common.FailWithMsg("无法读取文件")
		return
	}
	user := GetCurrentUserFromCache(c)
	username := user.(system.SysUser).Username
	fileName := username + "_avatar" + path.Ext(file.Filename)
	imgPath := common.Conf.Upload.SaveDir + "/avatar/" + fileName
	err = c.SaveUploadedFile(file, imgPath)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 将头像url保存到数据库
	query := common.Mysql.Where("username = ?", username).First(&user)
	err = query.Update("avatar", "/"+imgPath).Error
	if err != nil {
		common.FailWithMsg(err.Error())
	}
	resp := map[string]string{
		"name": fileName,
		"url":  "/" + imgPath,
	}

	common.SuccessWithData(resp)
}
