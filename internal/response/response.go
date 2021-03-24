package response

import (
	"go-xops/pkg/common"

	"github.com/gin-gonic/gin"
)

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

// PageInfo ...分页封装
type PageInfo struct {
	Current  uint  `json:"current" form:"current"`   // 当前页码
	PageSize uint  `json:"pageSize" form:"pageSize"` // 每页显示条数
	Total    int64 `json:"total"`                    // 数据总条数
	All      bool  `json:"all" form:"all"`           // 不使用分页
}

// PageData ....带分页数据封装
type PageData struct {
	PageInfo
	DataList interface{} `json:"data"` // 数据列表
}

// GetLimit ...计算limit/offset, 如果需要用到返回的PageSize, PageNum, 务必保证Total值有效
func (s *PageInfo) GetLimit() (int, int) {
	// 传入参数可能不合法, 设置默认值
	var pageSize int64
	var current int64
	total := s.Total
	// 每页显示条数不能小于1
	if s.PageSize < 1 {
		pageSize = 10
	} else {
		pageSize = int64(s.PageSize)
	}
	// 页码不能小于1
	if s.Current < 1 {
		current = 1
	} else {
		current = int64(s.Current)
	}

	// 如果偏移量比总条数还多
	if total > 0 && current > total {
		current = total
	}

	// 计算最大页码
	maxPageNum := total/pageSize + 1
	if total%pageSize == 0 {
		maxPageNum = total / pageSize
	}
	// 页码不能小于1
	if maxPageNum < 1 {
		maxPageNum = 1
	}

	// 超出最后一页
	if current > maxPageNum {
		current = maxPageNum
	}

	limit := pageSize
	offset := limit * (current - 1)

	// gorm v2参数从interface改为int, 这里也需要相应改变
	return int(limit), int(offset)
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
	c.Set(common.Conf.System.OperationLogKey, resp)
}
