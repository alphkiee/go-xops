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

// GetDepts doc
// @Summary Get /api/v1/dept/list
// @Description 查询所有部门
// @Produce json
// @Param data body request.DeptListReq true "name, creator, status"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/dept/list [get]
func GetDepts(c *gin.Context) {
	// 绑定参数
	var req request.DeptListReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	// 创建服务
	s := service.New()
	depts := s.GetDepts(&req)
	if req.Name != "" || req.Status != nil {
		var newResp []response.DictTreeResp
		utils.Struct2StructByJson(depts, &newResp)
		response.SuccessWithData(newResp)
	} else {
		var resp []response.DeptTreeResp
		resp = service.GenDeptTree(nil, depts)
		response.SuccessWithData(resp)
	}
}

// CreateDept doc
// @Summary Get /api/v1/dept/create
// @Description 创建部门
// @Produce json
// @Param data body request.CreateDeptReq true "name, sort, parent_id, creator"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/dept/create [post]
func CreateDept(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateDeptReq
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
	err = s.CreateDept(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateDeptById doc
// @Summary Get /api/v1/dept/update/:deptId
// @Description 更新部门
// @Produce json
// @Param data body request.UpdateDeptReq true "name, status, sort, parent_id"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/dept/update/:deptId [patch]
func UpdateDeptById(c *gin.Context) {
	// 绑定参数
	var req request.UpdateDeptReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	deptId := utils.Str2Uint(c.Param("deptId"))
	if deptId == 0 {
		response.FailWithMsg("部门编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateDeptById(deptId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteDeptByIds doc
// @Summary Delete /api/v1/dept/delete
// @Description 根据ID批量删除部门
// @Produce json
// @Param data body request.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/dept/delete [delete]
func BatchDeleteDeptByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteDeptByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
