package cmdb

import (
	s "go-xops/api/v1/system"
	"go-xops/assets/system"
	"go-xops/internal/service/cmd"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

// 获取列表
func GetHosts(c *gin.Context) {
	// 绑定参数
	var req cmd.HostListReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	hosts, err := cmd.GetHosts(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 转为commonStruct, 隐藏部分字段
	var respStruct []cmd.HostRep
	utils.Struct2StructByJson(hosts, &respStruct)
	// 返回分页数据
	var resp common.PageData
	// 设置分页参数
	resp.PageInfo = req.PageInfo
	// 设置数据列表
	resp.DataList = respStruct
	common.SuccessWithData(resp)
}

// 创建
func CreateHost(c *gin.Context) {
	user := s.GetCurrentUserFromCache(c)
	// 绑定参数
	var req cmd.CreateHostReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	m := make(map[string]string, 0)
	m["IP"] = "IP"
	m["AuthType"] = "认证类型"
	err = common.NewValidatorError(common.Validate.Struct(req), m)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 记录当前创建人信息
	req.Creator = user.(system.SysUser).Name
	err = cmd.CreateHost(&req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// 获取当前主机信息
func GetHostInfo(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	hostId := utils.Str2Uint(c.Param("id"))
	if hostId == 0 {
		common.FailWithMsg("接口编号不正确")
		return
	}
	host, err := cmd.GetHostByid(hostId)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	// 转为commonStruct, 隐藏部分字段
	var connStruct cmd.HostRep
	utils.Struct2StructByJson(host, &connStruct)
	common.SuccessWithData(connStruct)
}

// 更新
func UpdateHostById(c *gin.Context) {
	// 绑定参数
	var req gin.H
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	hostId := utils.Str2Uint(c.Param("id"))
	if hostId == 0 {
		common.FailWithMsg("接口编号不正确")
		return
	}
	err = cmd.UpdateHostById(hostId, req)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// 批删除
func BatchDeleteHostByIds(c *gin.Context) {
	var req common.IdsReq
	err := c.Bind(&req)
	if err != nil {
		common.FailWithCode(common.ParmError)
		return
	}
	err = cmd.DeleteHostsById(req.Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}

// ExcelIn 导入主机列表
func ExcelIn(c *gin.Context) {
	user := s.GetCurrentUserFromCache(c)
	u := user.(system.SysUser).Name
	dir := "/Users/痞老板/go/code/go-xops/upload/host/"
	//获取文件头
	file, err := c.FormFile("host")
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	//获取文件名
	fileName := file.Filename
	//保存文件到服务器本地
	if err := c.SaveUploadedFile(file, dir+fileName); err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	f := dir + fileName
	xlsx, err := excelize.OpenFile(f)
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}

	// 获取excel中具体的列的值
	rows, err := xlsx.GetRows("tb_cmdb_host")
	if err != nil {
		common.FailWithMsg(err.Error())
	}
	for key, row := range rows {
		if key > 0 {
			host := cmd.CreateHostReq{HostName: row[0], IP: row[1], HostType: row[4], Port: row[2], AuthType: row[5], User: row[6], Password: row[7], OsVersion: row[3], PrivateKey: row[8], Creator: u}
			err := cmd.CreateHost(&host)
			if err != nil {
				common.FailWithMsg(err.Error())
				return
			}
		}

	}
	common.Success()
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
		common.FailWithCode(common.ParmError)
		return
	}
	Ids := utils.Str2UintArr(c.Param("ids"))
	hosts, err := cmd.GetHostByIds(Ids)
	if err != nil {
		common.FailWithMsg(err.Error())
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
		common.FailWithMsg(err.Error())
	}

	titleRow := xlsx.AddRow()

	xlsRow := XlsxRow{
		Row:  titleRow,
		Data: t,
	}

	if xlsRow.Row == nil {
		common.FailWithMsg("xlsRow 无数据")
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
			common.FailWithMsg("xlsRow 数据为空")
			return
		}
		for _, v := range xlsRow.Data {
			cell := xlsRow.Row.AddCell()
			cell.SetString(v)
		}
	}

	err = file.Save("/Users/痞老板/go/code/go-xops/upload/host/cmdb_host.xlsx")
	if err != nil {
		common.FailWithMsg(err.Error())
		return
	}
	common.Success()
}
