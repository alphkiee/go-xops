package common

import "github.com/gin-gonic/gin"

// 自定义错误码与错误信息

const (
	// Ok ...
	Ok = 200
	// NotOk ...
	NotOk = 400
	// Unauthorized ...
	Unauthorized = 401
	// Forbidden ...
	Forbidden = 403
	// ParmError ...
	ParmError = 406
	// InternalServerError ...
	InternalServerError = 500
	// AuthError ...
	AuthError = 1000
	// UserForbidden ...
	UserForbidden = 1001
)

const (
	// OkMsg ...
	OkMsg = "操作成功"
	// NotOkMsg ...
	NotOkMsg = "操作失败"
	// UnauthorizedMsg ...
	UnauthorizedMsg = "登录过期, 需要重新登录"
	// LoginCheckErrorMsg ...
	LoginCheckErrorMsg = "用户名或密码错误"
	// ForbiddenMsg ...
	ForbiddenMsg = "无权访问该资源"
	// InternalServerErrorMsg ...
	InternalServerErrorMsg = "服务器内部错误"
	// ParmErrorMsg ...
	ParmErrorMsg = "参数绑定失败, 请检查数据类型"
	// UserForbiddenMsg ...
	UserForbiddenMsg = "用户已被禁用"
)

// CustomError ...
var CustomError = map[int]string{
	Ok:                  OkMsg,
	NotOk:               NotOkMsg,
	Unauthorized:        UnauthorizedMsg,
	Forbidden:           ForbiddenMsg,
	InternalServerError: InternalServerErrorMsg,
	AuthError:           LoginCheckErrorMsg,
	ParmError:           ParmErrorMsg,
	UserForbidden:       UserForbiddenMsg,
}

// 适用于前端传过来的
type IdsReq struct {
	Ids []uint `json:"ids" form:"ids"` // 传多个id
}

// 适用于前端传过来的
type KeyReq struct {
	Key string `json:"key" form:"key"` // 传多个id
}

// RespInfo ...http请求响应封装
type RespInfo struct {
	Code    int         `json:"code"`    // 错误代码代码
	Status  bool        `json:"status"`  // 状态
	Data    interface{} `json:"data"`    // 数据内容
	Message string      `json:"message"` // 消息提示
}

// RespPageInfo ...
type RespPageInfo struct {
	Code    int         `json:"code"`    // 错误代码代码
	Status  bool        `json:"status"`  // 状态
	Data    interface{} `json:"data"`    // 数据内容
	Message string      `json:"message"` // 消息提示
}

// Result ...
func Result(code int, status bool, data interface{}) {
	// 结果以panic异常的形式抛出, 交由异常处理中间件处理
	panic(RespInfo{
		Code:    code,
		Status:  status,
		Data:    data,
		Message: CustomError[code],
	})
}

// MsgResult ...
func MsgResult(code int, status bool, msg string, data interface{}) {
	// 结果以panic异常的形式抛出, 交由异常处理中间件处理
	panic(RespInfo{
		Code:    code,
		Status:  status,
		Data:    data,
		Message: msg,
	})
}

// Success ...
func Success() {
	Result(Ok, true, map[string]interface{}{})
}

// SuccessWithData ...
func SuccessWithData(data interface{}) {
	Result(Ok, true, data)
}

// SuccessWithPageData ...
func SuccessWithPageData(data interface{}) {
	Result(Ok, true, data)
}

// SuccessWithMsg ...
func SuccessWithMsg(msg string) {
	MsgResult(Ok, true, msg, map[string]interface{}{})
}

// SuccessWithCode ...
func SuccessWithCode(code int) {
	Result(code, true, map[string]interface{}{})
}

// FailWithMsg ...
func FailWithMsg(msg string) {
	MsgResult(NotOk, false, msg, map[string]interface{}{})
}

// FailWithCode ...
func FailWithCode(code int) {
	Result(code, false, map[string]interface{}{})
}

// JSON ...写入json返回值
func JSON(c *gin.Context, code int, resp interface{}) {
	// 调用gin写入json
	c.JSON(code, resp)
	// 保存响应对象到context, Operation Log会读取到
	c.Set(Conf.System.OperationLogKey, resp)
}
