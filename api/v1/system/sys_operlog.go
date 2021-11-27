package system

import (
	"go-xops/assets/system"
	s "go-xops/internal/service/system"
	"go-xops/pkg/cache"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetOperLogs doc
// @Summary Get /api/v1/operlog/list
// @Description 列出所有操作日志
// @Produce json
// @Param name query string false "name"
// @Param method query string false "method"
// @Param path query string false "path"
// @Param username query string false "username"
// @Param ip query string false "ip"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/operlog/list [get]
func GetOperLogs(c *gin.Context) {
	// 绑定参数
	var req s.OperLogListReq
	reqErr := c.Bind(&req)
	if reqErr != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	var operationLogs []system.SysOperLog
	var err error
	// 创建缓存对象
	cache, err := cache.New(time.Second * 20)
	if err != nil {
		logrus.Error(err)
	}
	key := "operationLog:" + req.Name + ":" + req.Method + ":" + req.Username + ":" + req.Ip + ":" + req.Path + ":" +
		strconv.Itoa(int(req.Current)) + ":" + strconv.Itoa(int(req.PageSize)) + ":" + strconv.Itoa(int(req.Total))

	cache.DBGetter = func() interface{} {
		operationLogs, err = s.GetOperLogs(&req)
		return operationLogs
	}
	// 获取缓存
	cache.Get(key)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 转为commonStruct, 隐藏部分字段
	var respStruct []s.OperationLogListResp
	utils.Struct2StructByJson(operationLogs, &respStruct)
	// 返回分页数据
	var resp common.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	common.SuccessWithData(resp)
}

// BatchDeleteOperLogByIds doc
// @Summary Delete /api/v1/operlog/delete
// @Description 根据ID批量删除日志
// @Produce json
// @Param data body common.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/operlog/delete [delete]
func BatchDeleteOperLogByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = s.DeleteOperationLogByIds(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}
