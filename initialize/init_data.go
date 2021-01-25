package initialize

import (
	"errors"
	"go-xops/dto/service"
	"go-xops/models"
	"go-xops/models/system"
	"go-xops/pkg/common"
	"go-xops/pkg/utils"

	"gorm.io/gorm"
)

// 初始化数据
func InitData() {
	// 1. 初始化角色
	creator := "系统创建"
	status := true
	roles := []system.SysRole{
		{
			Model: models.Model{
				Id: 1,
			},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    "系统管理员",
			Status:  &status,
			Creator: creator,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Name:    "访客",
			Keyword: "guest",
			Desc:    "外来访问人员",
			Status:  &status,
			Creator: creator,
		},
	}
	for _, role := range roles {
		oldRole := system.SysRole{}
		err := common.Mysql.Where("id = ?", role.Id).First(&oldRole).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.Mysql.Create(&role)
		}
	}

	// 2. 初始化菜单
	menus := []system.SysMenu{
		{
			Model: models.Model{
				Id: 1,
			},
			Name:     "仪表盘",
			Icon:     "iconyibiaopan",
			Path:     "index",
			Sort:     0,
			Status:   &status,
			ParentId: 0,
			Creator:  creator,
			Roles:    roles,
		},
		{
			Model: models.Model{
				Id: 10,
			},
			Name:     "资产管理",
			Icon:     "iconxitongshezhi1",
			Path:     "asset",
			Sort:     1,
			Status:   &status,
			ParentId: 0,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 11,
			},
			Name:     "主机管理",
			Path:     "host",
			Sort:     1,
			Status:   &status,
			ParentId: 10,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 12,
			},
			Name:     "连接管理",
			Path:     "connection",
			Sort:     2,
			Status:   &status,
			ParentId: 10,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Name:     "系统设置",
			Icon:     "iconxitongshezhi1",
			Path:     "system",
			Sort:     999,
			Status:   &status,
			ParentId: 0,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:     "用户管理",
			Path:     "user",
			Sort:     10,
			Status:   &status,
			ParentId: 2,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 4,
			},
			Name:     "部门管理",
			Path:     "dept",
			Sort:     11,
			Status:   &status,
			ParentId: 2,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 5,
			},
			Name:     "菜单管理",
			Path:     "menu",
			Sort:     12,
			Status:   &status,
			ParentId: 2,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 6,
			},
			Name:     "角色管理",
			Path:     "role",
			Sort:     13,
			Status:   &status,
			ParentId: 2,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 7,
			},
			Name:     "接口管理",
			Path:     "api",
			Sort:     14,
			Status:   &status,
			ParentId: 2,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 8,
			},
			Name:     "数据字典",
			Path:     "dict",
			Sort:     15,
			Status:   &status,
			ParentId: 2,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 9,
			},
			Name:     "操作日志",
			Path:     "operlog",
			Sort:     16,
			Status:   &status,
			ParentId: 2,
			Creator:  creator,
			Roles: []system.SysRole{
				roles[0],
			},
		},
	}
	for _, menu := range menus {
		oldMenu := system.SysMenu{}
		err := common.Mysql.Where("id = ?", menu.Id).First(&oldMenu).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.Mysql.Create(&menu)
		}
	}

	// 3. 初始化用户
	// 默认头像
	avatar := "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	users := []system.SysUser{
		{
			Username: "admin",
			Password: utils.GenPwd("123456"),
			Mobile:   "18888888888",
			Avatar:   avatar,
			Name:     "管理员",
			RoleId:   1,
			Creator:  creator,
		},
		{
			Username: "guest",
			Password: utils.GenPwd("123456"),
			Mobile:   "15888888888",
			Avatar:   avatar,
			Name:     "来宾",
			RoleId:   2,
			Creator:  creator,
		},
	}
	for _, user := range users {
		oldUser := system.SysUser{}
		err := common.Mysql.Where("username = ?", user.Username).First(&oldUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.Mysql.Create(&user)
		}
	}
	// 初始化字典
	dicts := []system.SysDict{
		{
			Model: models.Model{
				Id: 1,
			},
			Key:     "env_type",
			Value:   "应用环境",
			Desc:    "应用环境",
			Creator: creator,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Key:      "dev",
			Value:    "开发环境",
			ParentId: 1,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Key:      "test",
			Value:    "测试环境",
			ParentId: 1,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 4,
			},
			Key:      "prod",
			Value:    "生产环境",
			ParentId: 1,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 5,
			},
			Key:     "auth_type",
			Value:   "认证类型",
			Desc:    "主机认证类型",
			Creator: creator,
		},
		{
			Model: models.Model{
				Id: 6,
			},
			Key:      "key",
			Value:    "秘钥验证",
			ParentId: 5,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 7,
			},
			Key:      "password",
			Value:    "密码验证",
			ParentId: 5,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 8,
			},
			Key:     "host_type",
			Value:   "主机类型",
			Desc:    "主机分类",
			Creator: creator,
		},
		{
			Model: models.Model{
				Id: 9,
			},
			Key:      "vm",
			Value:    "虚拟机",
			ParentId: 8,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 9,
			},
			Key:      "phy",
			Value:    "物理机",
			ParentId: 8,
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 9,
			},
			Key:      "aliyun",
			Value:    "阿里云",
			ParentId: 8,
			Creator:  creator,
		},
	}

	for _, dict := range dicts {
		oldDict := system.SysDict{}
		err := common.Mysql.Where("id = ?", dict.Id).First(&oldDict).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.Mysql.Create(&dict)
		}
	}

	// 5. 初始化接口
	apis := []system.SysApi{
		{
			Model: models.Model{
				Id: 1,
			},
			Name:     "用户登录",
			Method:   "POST",
			Path:     "/auth/login",
			Category: "基本权限",
			Desc:     "获取用户信息和token",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 2,
			},
			Name:     "用户登出",
			Method:   "POST",
			Path:     "/auth/logout",
			Category: "基本权限",
			Desc:     "用户登出",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 3,
			},
			Name:     "刷新令牌",
			Method:   "POST",
			Path:     "/auth/refresh_token",
			Category: "基本权限",
			Desc:     "刷新JWT令牌",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 4,
			},
			Name:     "用户信息",
			Method:   "GET",
			Path:     "/v1/user/info",
			Category: "基本权限",
			Desc:     "用户信息",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 5,
			},
			Name:     "更新基本信息",
			Method:   "PATCH",
			Path:     "/v1/user/info/update/:userId",
			Category: "基本权限",
			Desc:     "更新信息",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 6,
			},
			Name:     "上传头像",
			Method:   "POST",
			Path:     "/v1/user/info/uploadImg",
			Category: "基本权限",
			Desc:     "上传头像",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 7,
			},
			Name:     "修改密码",
			Method:   "PUT",
			Path:     "/v1/user/changePwd",
			Category: "基本权限",
			Desc:     "修改密码",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 8,
			},
			Name:     "查询用户",
			Method:   "GET",
			Path:     "/v1/user/list",
			Category: "用户管理",
			Desc:     "用户列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 9,
			},
			Name:     "创建用户",
			Method:   "POST",
			Path:     "/v1/user/create",
			Category: "用户管理",
			Desc:     "创建用户",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 10,
			},
			Name:     "更新用户",
			Method:   "PATCH",
			Path:     "/v1/user/update/:userId",
			Category: "用户管理",
			Desc:     "更新用户",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 11,
			},
			Name:     "删除用户",
			Method:   "DELETE",
			Path:     "/v1/user/delete",
			Category: "用户管理",
			Desc:     "删除用户",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 12,
			},
			Name:     "当前菜单",
			Method:   "GET",
			Path:     "/v1/menu/tree",
			Category: "基本权限",
			Desc:     "获取菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 13,
			},
			Name:     "查询菜单",
			Method:   "GET",
			Path:     "/v1/menu/list",
			Category: "菜单管理",
			Desc:     "菜单列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 14,
			},
			Name:     "创建菜单",
			Method:   "POST",
			Path:     "/v1/menu/create",
			Category: "菜单管理",
			Desc:     "创建菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 15,
			},
			Name:     "更新菜单",
			Method:   "PATCH",
			Path:     "/v1/menu/update/:menuId",
			Category: "菜单管理",
			Desc:     "更新菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 16,
			},
			Name:     "删除菜单",
			Method:   "DELETE",
			Path:     "/v1/menu/delete",
			Category: "菜单管理",
			Desc:     "删除菜单",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 17,
			},
			Name:     "查询角色",
			Method:   "GET",
			Path:     "/v1/role/list",
			Category: "角色管理",
			Desc:     "角色列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 18,
			},
			Name:     "创建角色",
			Method:   "POST",
			Path:     "/v1/role/create",
			Category: "角色管理",
			Desc:     "创建角色",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 19,
			},
			Name:     "修改角色",
			Method:   "PATCH",
			Path:     "/v1/role/update/:roleId",
			Category: "角色管理",
			Desc:     "更新角色",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 20,
			},
			Name:     "修改权限",
			Method:   "PATCH",
			Path:     "/v1/role/perms/update/:roleId",
			Category: "角色管理",
			Desc:     "更新权限",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 21,
			},
			Name:     "获取权限",
			Method:   "GET",
			Path:     "/v1/role/perms/:roleId",
			Category: "角色管理",
			Desc:     "获取权限信息",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 22,
			},
			Name:     "删除角色",
			Method:   "DELETE",
			Path:     "/v1/role/delete",
			Category: "角色管理",
			Desc:     "删除角色",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 23,
			},
			Name:     "查询接口",
			Method:   "GET",
			Path:     "/v1/api/list",
			Category: "接口管理",
			Desc:     "获取接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 24,
			},
			Name:     "创建接口",
			Method:   "POST",
			Path:     "/v1/api/create",
			Category: "接口管理",
			Desc:     "创建接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 25,
			},
			Name:     "修改接口",
			Method:   "PATCH",
			Path:     "/v1/api/update/:apiId",
			Category: "接口管理",
			Desc:     "更新接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 26,
			},
			Name:     "删除接口",
			Method:   "DELETE",
			Path:     "/v1/api/delete",
			Category: "接口管理",
			Desc:     "删除接口",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 27,
			},
			Name:     "查询日志",
			Method:   "GET",
			Path:     "/v1/operlog/list",
			Category: "日志管理",
			Desc:     "获取操作日志",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 28,
			},
			Name:     "删除日志",
			Method:   "DELETE",
			Path:     "/v1/operlog/delete",
			Category: "日志管理",
			Desc:     "删除操作日志",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 29,
			},
			Name:     "查询字典",
			Method:   "GET",
			Path:     "/v1/dict/list",
			Category: "字典管理",
			Desc:     "字典列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 30,
			},
			Name:     "创建字典",
			Method:   "POST",
			Path:     "/v1/dict/create",
			Category: "字典管理",
			Desc:     "创建字典",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 31,
			},
			Name:     "更新字典",
			Method:   "PATCH",
			Path:     "/v1/dict/update/:dictId",
			Category: "字典管理",
			Desc:     "更新字典",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 32,
			},
			Name:     "删除字典",
			Method:   "DELETE",
			Path:     "/v1/dict/delete",
			Category: "字典管理",
			Desc:     "删除字典",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 33,
			},
			Name:     "查询主机",
			Method:   "GET",
			Path:     "/v1/host/list",
			Category: "主机管理",
			Desc:     "主机列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 34,
			},
			Name:     "创建主机",
			Method:   "POST",
			Path:     "/v1/host/create",
			Category: "主机管理",
			Desc:     "创建主机",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 35,
			},
			Name:     "更新主机",
			Method:   "PATCH",
			Path:     "/v1/host/update/:hostId",
			Category: "主机管理",
			Desc:     "更新主机",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 36,
			},
			Name:     "删除主机",
			Method:   "DELETE",
			Path:     "/v1/host/delete",
			Category: "主机管理",
			Desc:     "删除主机",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 37,
			},
			Name:     "查询连接",
			Method:   "GET",
			Path:     "/v1/host/connection/list",
			Category: "连接管理",
			Desc:     "连接列表",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 38,
			},
			Name:     "注销连接",
			Method:   "DELETE",
			Path:     "/v1/host/connection/delete",
			Category: "连接管理",
			Desc:     "注销连接",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 39,
			},
			Name:     "连接SSH",
			Method:   "GET",
			Path:     "/v1/host/ssh",
			Category: "主机管理",
			Desc:     "连接ssh",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 40,
			},
			Name:     "显示文件",
			Method:   "GET",
			Path:     "/v1/host/ssh/ls",
			Category: "文件管理",
			Desc:     "显示文件",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 41,
			},
			Name:     "上传文件",
			Method:   "POST",
			Path:     "/v1/host/ssh/upload",
			Category: "文件管理",
			Desc:     "上传文件",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 42,
			},
			Name:     "下载文件",
			Method:   "GET",
			Path:     "/v1/host/ssh/download",
			Category: "文件管理",
			Desc:     "下载文件",
			Creator:  creator,
		},
		{
			Model: models.Model{
				Id: 43,
			},
			Name:     "删除文件",
			Method:   "DELETE",
			Path:     "/v1/host/ssh/rm",
			Category: "文件管理",
			Desc:     "删除文件",
			Creator:  creator,
		},
	}
	for _, api := range apis {
		oldApi := system.SysApi{}
		err := common.Mysql.Where("id = ?", api.Id).First(&oldApi).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.Mysql.Create(&api)
			s := service.New()
			// 管理员拥有所有API权限role[0]
			//_, err = s.CreateRoleCasbin(system.SysRoleCasbin{
			//	Keyword: roles[0].Keyword,
			//	Path:    api.Path,
			//	Method:  api.Method,
			//})
			// 来宾权限
			if api.Id <= 7 || api.Id == 12 {
				_, err = s.CreateRoleCasbin(system.SysRoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
		}
	}
}
