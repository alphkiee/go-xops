package system

import (
	"go-xops/assets/system"
	s "go-xops/internal/service/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetDepts doc
// @Summary Get /api/v1/dept/list
// @Description 查询所有部门
// @Produce json
// @Param data body s.DeptListReq true "name, creator, status"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/dept/list [get]
func GetDepts(c *gin.Context) {
	// 绑定参数
	var req s.DeptListReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	depts := s.GetDepts(&req)
	if req.Name != "" || req.Status != nil {
		var newResp []s.DictTreeResp
		utils.Struct2StructByJson(depts, &newResp)
		common.SuccessWithData(newResp)
	} else {
		var resp []s.DeptTreeResp
		resp = s.GenDeptTree(nil, depts)
		common.SuccessWithData(resp)
	}
}

// CreateDept doc
// @Summary Get /api/v1/dept/create
// @Description 创建部门
// @Produce json
// @Param data body s.CreateDeptReq true "name, sort, parent_id, creator"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/dept/create [post]
func CreateDept(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req s.CreateDeptReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["Name"] = "姓名"
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	err = s.CreateDept(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateDeptById doc
// @Summary Get /api/v1/dept/update/:deptId
// @Description 更新部门
// @Produce json
// @Param data body s.UpdateDeptReq true "name, status, sort, parent_id"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/dept/update/:deptId [patch]
func UpdateDeptById(c *gin.Context) {
	// 绑定参数
	var req s.UpdateDeptReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	deptId := utils.Str2Uint(c.Param("deptId"))
	if deptId == 0 {
		common.FailWithMsg("部门编号不正确")
		return
	}
	err = s.UpdateDeptById(deptId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// BatchDeleteDeptByIds doc
// @Summary Delete /api/v1/dept/delete
// @Description 根据ID批量删除部门
// @Produce json
// @Param data body common.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/dept/delete [delete]
func BatchDeleteDeptByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = s.DeleteDeptByIds(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}
