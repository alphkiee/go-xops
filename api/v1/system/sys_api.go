package system

import (
	"go-xops/assets/system"
	s "go-xops/internal/service/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetApis doc
// @Summary Get /api/v1/api/list
// @Description 查看所有API
// @Produce json
// @Param data body s.ApiListReq true "name, method, path, category, creator, tree"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/api/list [post]
func GetApis(c *gin.Context) {
	// 绑定参数
	var req s.ApiListReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}

	apis, err := s.GetApis(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	var respStruct []s.ApiListResp
	utils.Struct2StructByJson(apis, &respStruct)
	if req.Tree {
		// 转换成树结构
		tree := make([]s.ApiTreeResp, 0)
		for _, api := range respStruct {
			existIndex := -1
			children := make([]s.ApiListResp, 0)
			for index, leaf := range tree {
				if leaf.Category == api.Category {
					children = leaf.Children
					existIndex = index
					break
				}
			}
			// api结构转换
			var item s.ApiListResp
			utils.Struct2StructByJson(api, &item)
			children = append(children, item)
			if existIndex != -1 {
				// 更新元素
				tree[existIndex].Children = children
			} else {
				// 新增元素
				tree = append(tree, s.ApiTreeResp{
					Category: api.Category,
					Children: children,
				})
			}
		}

		common.SuccessWithData(tree)
		return
	}
	// 返回分页数据
	var resp common.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	common.SuccessWithData(resp)
}

// CreateApi doc
// @Summary Get /api/v1/api/create
// @Description 创建api
// @Produce json
// @Param data body s.CreateApiReq true "name, method, path, category, creator, desc"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/api/create [post]
func CreateApi(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req s.CreateApiReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["Name"] = "姓名"
	m["Method"] = "方法"
	m["Path"] = "路径"
	m["Category"] = "分类"
	// 参数校验
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	err = s.CreateApi(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// UpdateApiById doc
// @Summary Get /api/v1/api/update/:apiId
// @Description 更新api
// @Produce json
// @Param apiId path int true "apiId"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/api/update/:apiId [patch]
func UpdateApiById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}

	// 获取path中的apiId
	apiId := utils.Str2Uint(c.Param("apiId"))
	if apiId == 0 {
		common.FailWithMsg("接口编号不正确")
		return
	}
	err = s.UpdateApiById(apiId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// BatchDeleteApiByIds doc
// @Summary Delete /api/v1/api/delete
// @Description 根据ID批量删除api
// @Produce json
// @Param data body common.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} common.RespInfo
// @Failure 400 {object} common.RespInfo
// @Router /api/v1/api/delete [delete]
func BatchDeleteApiByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = s.DeleteApiByIds(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}
