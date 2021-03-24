package cmdb

import (
	"go-xops/api/v1/system"
	s "go-xops/assets/system"
	"go-xops/internal/request"
	"go-xops/internal/response"
	"go-xops/internal/service"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

// 获取列表
func GetHosts(c *gin.Context) {
	// 绑定参数
	var req request.HostListReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}

	// 创建服务
	s := service.New()
	hosts, err := s.GetHosts(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var respStruct []response.HostRep
	utils.Struct2StructByJson(hosts, &respStruct)
	// 返回分页数据
	var resp response.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	response.SuccessWithData(resp)
}

// 创建
func CreateHost(c *gin.Context) {
	user := system.GetCurrentUserFromCache(c)
	// 绑定参数
	var req request.CreateHostReq
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
	req.Creator = user.(s.SysUser).Name
	// 创建服务
	s := service.New()
	err = s.CreateHost(&req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 获取当前主机信息
func GetHostInfo(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	hostId := utils.Str2Uint(c.Param("id"))
	if hostId == 0 {
		response.FailWithMsg("接口编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	host, err := s.GetHostByid(hostId)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	// 转为ResponseStruct, 隐藏部分字段
	var connStruct response.HostRep
	utils.Struct2StructByJson(host, &connStruct)
	response.SuccessWithData(connStruct)
}

// 更新
func UpdateHostById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	hostId := utils.Str2Uint(c.Param("id"))
	if hostId == 0 {
		response.FailWithMsg("接口编号不正确")
		return
	}
	// 创建服务
	s := service.New()
	// 更新数据
	err = s.UpdateHostById(hostId, req)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// 批删除
func BatchDeleteHostByIds(c *gin.Context) {
	var req request.IdsReq
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	// 创建服务
	s := service.New()
	// 删除数据
	err = s.DeleteHostsById(req.Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}

// ExcelIn 导入主机列表
func ExcelIn(c *gin.Context) {
	user := system.GetCurrentUserFromCache(c)
	u := user.(s.SysUser).Name
	dir := "/Users/痞老板/go/code/go-xops/upload/host/"
	//获取文件头
	file, err := c.FormFile("host")
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	//获取文件名
	fileName := file.Filename
	//保存文件到服务器本地
	if err := c.SaveUploadedFile(file, dir+fileName); err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	f := dir + fileName
	xlsx, err := excelize.OpenFile(f)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// 获取excel中具体的列的值
	rows, err := xlsx.GetRows("tb_cmdb_host")
	if err != nil {
		response.FailWithMsg(err.Error())
	}
	for key, row := range rows {
		if key > 0 {
			host := request.CreateHostReq{HostName: row[0], IP: row[1], HostType: row[4], Port: row[2], AuthType: row[5], User: row[6], Password: row[7], OsVersion: row[3], PrivateKey: row[8], Creator: u}

			s := service.New()
			err := s.CreateHost(&host)
			if err != nil {
				response.FailWithMsg(err.Error())
				return
			}
		}

	}
	response.Success()
}

// ExportHost 导出主机列表
func ExportHost(c *gin.Context) {

	type XlsxRow struct {
		Row  *xlsx.Row
		Data []string
	}

	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		response.FailWithCode(response.ParmError)
		return
	}
	Ids := utils.Str2UintArr(c.Param("ids"))

	s := service.New()
	hosts, err := s.GetHostByIds(Ids)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	t := make([]string, 0)
	t = append(t, "host_name")
	t = append(t, "ip")
	t = append(t, "os_version")
	t = append(t, "auth_type")
	t = append(t, "creator")

	file := xlsx.NewFile()
	xlsx, err := file.AddSheet("Sheet")
	if err != nil {
		response.FailWithMsg(err.Error())
	}

	titleRow := xlsx.AddRow()

	xlsRow := XlsxRow{
		Row:  titleRow,
		Data: t,
	}

	if xlsRow.Row == nil {
		response.FailWithMsg("xlsRow 无数据")
		return
	}
	for _, v := range xlsRow.Data {
		cell := xlsRow.Row.AddCell()
		cell.SetString(v)
	}
	for _, v := range hosts {
		currentRow := xlsx.AddRow()
		tmp := make([]string, 0)
		tmp = append(tmp, v.HostName)
		tmp = append(tmp, v.IP)
		tmp = append(tmp, v.OsVersion)
		tmp = append(tmp, v.AuthType)
		tmp = append(tmp, v.Creator)

		xlsRow := XlsxRow{
			Row:  currentRow,
			Data: tmp,
		}
		if xlsRow.Data == nil {
			response.FailWithMsg("xlsRow 数据为空")
			return
		}
		for _, v := range xlsRow.Data {
			cell := xlsRow.Row.AddCell()
			cell.SetString(v)
		}
	}

	err = file.Save("/Users/痞老板/go/code/go-xops/upload/host/cmdb_host.xlsx")
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	response.Success()
}
