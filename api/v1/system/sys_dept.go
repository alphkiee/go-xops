package system

import (
	"go-xops/dto/request"
	"go-xops/dto/response"
	"go-xops/dto/service"
	"go-xops/models/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// 查询所有部门
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

// 创建部门
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

// 更新部门
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

// 批量删除部门
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
