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

// GetApis doc
// @Summary Get /api/v1/api/list
// @Description 查看所有API
// @Produce json
// @Param data body request.ApiListReq true "name, method, path, category, creator, tree"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/api/list [post]
func GetApis(c *gin.Context) {
	// 绑定参数
	var req request.ApiListReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	apis, err := s.GetApis(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.ApiListResp
	utils.Struct2StructByJson(apis, &respStruct)
	if req.Tree {
		// 转换成树结构
		tree := make([]response.ApiTreeResp, 0)
		for _, api := range respStruct {
			existIndex := -1
			children := make([]response.ApiListResp, 0)
			for index, leaf := range tree {
				if leaf.Category == api.Category {
					children = leaf.Children
					existIndex = index
					break
				}
			}
			// api结构转换
			var item response.ApiListResp
			utils.Struct2StructByJson(api, &item)
			children = append(children, item)
			if existIndex != -1 {
				// 更新元素
				tree[existIndex].Children = children
			} else {
				// 新增元素
				tree = append(tree, response.ApiTreeResp{
					Category: api.Category,
					Children: children,
				})
			}
		}

		response.SuccessWithData(tree)
		return
	}
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response.SuccessWithData(resp)
}

// CreateApi doc
// @Summary Get /api/v1/api/create
// @Description 创建api
// @Produce json
// @Param data body request.CreateApiReq true "name, method, path, category, creator, desc"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/api/create [post]
func CreateApi(c *gin.Context) {
	user := GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateApiReq
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
	err = s.CreateApi(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// UpdateApiById doc
// @Summary Get /api/v1/api/update/:apiId
// @Description 更新api
// @Produce json
// @Param apiId path int true "apiId"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/api/update/:apiId [patch]
func UpdateApiById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 获取path中的apiId
	apiId := utils.Str2Uint(c.Param("apiId"))
	if apiId == 0 {
		response.FailWithMsg("接口编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateApiById(apiId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// BatchDeleteApiByIds doc
// @Summary Delete /api/v1/api/delete
// @Description 根据ID批量删除api
// @Produce json
// @Param data body request.IdsReq true "ids"
// @Security ApiKeyAuth
// @Success 200 {object} response.RespInfo
// @Failure 400 {object} response.RespInfo
// @Router /api/v1/api/delete [delete]
func BatchDeleteApiByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteApiByIds(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
