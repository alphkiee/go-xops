package initialize

import (
	"go-xops/pkg/common"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// Casbin ...获取casbin策略管理器
func Casbin() {
	// 初始化数据库适配器, 添加自定义表前缀, casbin不使用事务管理, 因为他内部使用到事务, 重复用会导致冲突
	// casbin默认表名casbin_rule, 为了与项目统一改写一下规则
	// 注意: gormadapter.CasbinTableName内部添加了下划线, 这里不再多此一举
	//a, err := gormadapter.NewAdapterByDBUseTableName(common.Mysql, common.Conf.Mysql.TablePrefix, "sys_casbin")
	a, err := gormadapter.NewAdapterByDB(common.Mysql)
	if err != nil {
		panic(err)
	}
	// 读取配置文件
	config, err := common.ConfBox.Find(common.Conf.Casbin.ModelPath)
	cabinModel := model.NewModel()
	// 从字符串中加载casbin配置
	err = cabinModel.LoadModelFromText(string(config))
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewSyncedEnforcer(cabinModel, a)
	if err != nil {
		panic(err)
	}
	// 加载策略
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}
	//if common.Conf.System.AppMode == "debug" {
	//	e.EnableLog(true)
	//}
	e.StartAutoLoadPolicy(time.Duration(common.Conf.Casbin.LoadDelay) * time.Second)
	common.Casbin = e
	// 关闭
	//defer common.Casbin.StopAutoLoadPolicy()
}
