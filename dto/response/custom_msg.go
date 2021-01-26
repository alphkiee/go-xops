package response

// 自定义错误码与错误信息

const (
	Ok                  = 200
	NotOk               = 400
	Unauthorized        = 401
	Forbidden           = 403
	ParmError           = 406
	InternalServerError = 500
	AuthError           = 1000
	UserForbidden = 1001
)

const (
	OkMsg                  = "操作成功"
	NotOkMsg               = "操作失败"
	UnauthorizedMsg        = "登录过期, 需要重新登录"
	LoginCheckErrorMsg     = "用户名或密码错误"
	ForbiddenMsg           = "无权访问该资源"
	InternalServerErrorMsg = "服务器内部错误"
	ParmErrorMsg           = "参数绑定失败, 请检查数据类型"
	UserForbiddenMsg          = "用户已被禁用"
)

var CustomError = map[int]string{
	Ok:                  OkMsg,
	NotOk:               NotOkMsg,
	Unauthorized:        UnauthorizedMsg,
	Forbidden:           ForbiddenMsg,
	InternalServerError: InternalServerErrorMsg,
	AuthError:           LoginCheckErrorMsg,
	ParmError:           ParmErrorMsg,
	UserForbidden:UserForbiddenMsg,
}
