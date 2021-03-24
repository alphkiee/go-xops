package system

import (
	"go-xops/assets/system"
	"go-xops/internal/request"
	"go-xops/internal/response"
	"go-xops/internal/service"
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
	u, _ := user.(response.LoginResp)
	// 创建缓存对象
	cache, err := cache.New(time.Second * 15)
	if err != nil {
		logrus.Error(err)
	}
	key := "user:" + u.Username
	cache.DBGetter = func() interface{} {
		// 创建mysql服务
		s := service.New()
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
// @Success 200 {object} response.RespInfo
// @Router /api/v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 转为UserInfoResponseStruct, 隐藏部分字段
	var resp response.UserInfoResp
	utils.Struct2StructByJson(user, &resp)
	response.SuccessWithData(resp)
}

// CreateUser doc
// @Summary Post /api/v1/user/create
// @Description 创建用户
// @Produce json
// @Param data body request.CreateUserReq true "username, password, name, role_id"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/user/create [post]
func CreateUser(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateUserReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 断言，创建结构体获取当前创建人信息
	req.Creator = user.(system.SysUser).Name
	// 创建服务
	s := service.New()
	err = s.CreateUser(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
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
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/user/list [get]
func GetUsers(c *gin.Context) {
	// 绑定参数
	var req request.UserListReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	users, err := s.GetUsers(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.UserListResp
	utils.Struct2StructByJson(users, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response.SuccessWithData(resp)

}

// UpdateUserBaseInfoById doc
// @Summary Patch /api/v1/user/info/update/:userId
// @Description 根据user ID来更新用户基本信息
// @Produce json
// @Param userId path int true "userId"
// @Param data body request.UpdateUserBaseInfoReq true "mobile, name, email"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/user/info/update/{userId} [patch]
func UpdateUserBaseInfoById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateUserBaseInfoReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 获取url path中的userId
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		response.FailWithMsg("用户编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateUserBaseInfoById(userId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateUserById doc
// @Summary Patch /api/v1/user/update/:userId
// @Description 更新用户根据 user ID
// @Produce json
// @Param userId path int true "userId"
// @Param data body request.UpdateUserReq true "mobile, name, email, password"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/user/update/{userId} [patch]
func UpdateUserById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateUserReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 获取url path中的userId
	userId := utils.Str2Uint(c.Param("userId"))
	if userId == 0 {
		response.FailWithMsg("用户编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateUserById(userId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// ChangePwd doc
// @Summary Put /api/v1/user/changePwd
// @Description 更改用户的密码
// @Produce json
// @Param data body request.ChangePwdReq true "old_password, new_password"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/user/changePwd [put]
func ChangePwd(c *gin.Context) {
	var msg string
	// 请求json绑定
	var req request.ChangePwdReq
	_ = c.ShouldBindJSON(&req)
	// 参数校验
	err := common.NewValidatorError(common.Validate.Struct(req), req.FieldTrans())
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 获取当前用户
	user := GetCurrentUserFromCache(c)
	query := common.Mysql.Where("username = ?", user.(system.SysUser).Username).First(&user)
	// 查询用户
	err = query.Error
	if err != nil {
		msg = err.Error()
	} else {
		// 校验密码
		if ok := utils.ComparePwd(req.OldPassword, user.(system.SysUser).Password); !ok {
			msg = "原密码错误"
		} else {
			// 更新密码
			err = query.Update("password", utils.GenPwd(req.NewPassword)).Error
			if err != nil {
				msg = err.Error()
			}
		}
	}
	if msg != "" {
		response.FailWithMsg(msg)
		return
	}
	response.Success()
}

// DeleteUserByIds doc
// @Summary Delete /api/v1/user/delete
// @Description 根据ID批量删除用户
// @Produce json
// @Param data body request.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/user/delete [delete]
func DeleteUserByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteUserByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UserAvatarUpload doc
// @Summary Post /api/v1/user/info/uploadImg
// @Description 上传头像
// @Produce json
// @Param avatar formData file true "avatar"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/user/info/uploadImg [post]
func UserAvatarUpload(c *gin.Context) {
	// 限制头像2MB(二进制移位xxxMB)
	err := c.Request.ParseMultipartForm(2 << 20)
	if err != nil {
		response.FailWithMsg("文件为空或图片大小超出最大值2MB")
		return
	}
	// 读取文件
	file, err := c.FormFile("avatar")
	if err != nil {
		response.FailWithMsg("无法读取文件")
		return
	}
	user := GetCurrentUserFromCache(c)
	username := user.(system.SysUser).Username
	fileName := username + "_avatar" + path.Ext(file.Filename)
	imgPath := common.Conf.Upload.SaveDir + "/avatar/" + fileName
	err = c.SaveUploadedFile(file, imgPath)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 将头像url保存到数据库
	query := common.Mysql.Where("username = ?", username).First(&user)
	err = query.Update("avatar", "/"+imgPath).Error
	if err != nil {
		response.FailWithMsg(err.Error())
	}
	resp := map[string]string{
		"name": fileName,
		"url":  "/" + imgPath,
	}

	response.SuccessWithData(resp)
}
