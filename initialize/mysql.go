package initialize

import (
	"fmt"
	"go-xops/assets/cmdb"
	"go-xops/assets/system"
	"go-xops/pkg/common"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化mysql数据库
func Mysql() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s&charset=%s&collation=%s",
		common.Conf.Mysql.Username,
		common.Conf.Mysql.Password,
		common.Conf.Mysql.Host,
		common.Conf.Mysql.Port,
		common.Conf.Mysql.Database,
		common.Conf.Mysql.Query,
		common.Conf.Mysql.Charset,
		common.Conf.Mysql.Collation,
	)
	if common.Conf.Mysql.LogMode {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			panic(fmt.Sprintf("初始化mysql异常: %v", err))
		}
		common.Mysql = db
	} else {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			panic(fmt.Sprintf("初始化mysql异常: %v", err))
		}
		common.Mysql = db
	}

	// 表结构
	autoMigrate()
}

// 自动迁移表结构
func autoMigrate() {
	common.Mysql.AutoMigrate(
		//new(models.SysCasbin),
		new(system.SysUser),
		new(system.SysDept),
		new(system.SysRole),
		new(system.SysMenu),
		new(system.SysDict),
		new(system.SysApi),
		new(system.SysOperLog),
		new(cmdb.Host),
	)
}
