package system

import (
	"go-xops/assets/system"
	s "go-xops/internal/service/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetDicts doc
// @Summary Get /api/v1/dict/list
// @Description 查询所有字典
// @Produce json
// @Param data body s.DictListReq true "key, value, desc, creator, status, type_key"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Router /api/v1/dict/list [get]
func GetDicts(c *gin.Context) {
	// 绑定参数
	var req s.DictListReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	dicts := s.GetDicts(&req)
	if req.Key != "" || req.Value != "" || req.Status != nil || req.TypeKey != "" {
		var newResp []s.DictTreeResp
		utils.Struct2StructByJson(dicts, &newResp)
		common.SuccessWithData(newResp)
	} else {
		var resp []s.DictTreeResp
		resp = s.GenDictTree(nil, dicts)
		common.SuccessWithData(resp)
	}
}

// CreateDict doc
// @Summary Get /api/v1/dict/create
// @Description 创建菜单
// @Produce json
// @Param data body s.CreateDictReq true "key, value, desc, parent_id, creator"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/dict/create [post]
func CreateDict(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req s.CreateDictReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["Key"] = "键"
	m["Value"] = "值"
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	err = s.CreateDict(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateDictById doc
// @Summary Get /api/v1/dict/update/:dictId
// @Description 更新字典
// @Produce json
// @Param data body s.UpdateDictReq true "key, value, desc, parent_id, status"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/dict/update/:dictId [patch]
func UpdateDictById(c *gin.Context) {
	// 绑定参数
	var req s.UpdateDictReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	dictId := utils.Str2Uint(c.Param("dictId"))
	if dictId == 0 {
		common.FailWithMsg("字典编号不正确")
		return
	}
	err = s.UpdateDictById(dictId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// BatchDeleteDictByIds doc
// @Summary Delete /api/v1/dict/delete
// @Description 根据ID批量删除菜单
// @Produce json
// @Param data body common.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/dict/delete [delete]
func BatchDeleteDictByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = s.DeleteDictByIds(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}
