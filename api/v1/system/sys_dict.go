package system

import (
	"go-xops/assets/system"
	"go-xops/internal/request"
	"go-xops/internal/response"
	"go-xops/internal/service"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetDicts doc
// @Summary Get /api/v1/dict/list
// @Description 查询所有字典
// @Produce json
// @Param data body request.DictListReq true "key, value, desc, creator, status, type_key"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Router /api/v1/dict/list [get]
func GetDicts(c *gin.Context) {
	// 绑定参数
	var req request.DictListReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	// 创建服务
	s := service.New()
	dicts := s.GetDicts(&req)
	if req.Key != "" || req.Value != "" || req.Status != nil || req.TypeKey != "" {
		var newResp []response.DictTreeResp
		utils.Struct2StructByJson(dicts, &newResp)
		response.SuccessWithData(newResp)
	} else {
		var resp []response.DictTreeResp
		resp = service.GenDictTree(nil, dicts)
		response.SuccessWithData(resp)
	}
}

// CreateDict doc
// @Summary Get /api/v1/dict/create
// @Description 创建菜单
// @Produce json
// @Param data body request.CreateDictReq true "key, value, desc, parent_id, creator"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/dict/create [post]
func CreateDict(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateDictReq
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
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	// 创建服务
	s := service.New()
	err = s.CreateDict(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateDictById doc
// @Summary Get /api/v1/dict/update/:dictId
// @Description 更新字典
// @Produce json
// @Param data body request.UpdateDictReq true "key, value, desc, parent_id, status"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/dict/update/:dictId [patch]
func UpdateDictById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateDictReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	dictId := utils.Str2Uint(c.Param("dictId"))
	if dictId == 0 {
		response.FailWithMsg("字典编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateDictById(dictId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteDictByIds doc
// @Summary Delete /api/v1/dict/delete
// @Description 根据ID批量删除菜单
// @Produce json
// @Param data body request.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/dict/delete [delete]
func BatchDeleteDictByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteDictByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
