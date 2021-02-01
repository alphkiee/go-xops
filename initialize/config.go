package initialize

import (
	"bytes"
	"fmt"
	"go-xops/pkg/common"
	"os"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

const (
	configBoxName = "xops-conf-box"
	configType    = "yml"
	// 配置文件目录, packr.Box基于当前包目录, 文件名需要写完整, 即使viper可以自动获取
	configPath = "../conf"
	devConfig  = "config-dev.yml"
	testConfig = "config-test.yml"
	prodConfig = "config-prod.yml"
)

// InitConfig ...初始化配置文件
func InitConfig() {
	// 使用packr将配置文件打包到二进制文件中, 如果以docker镜像方式运行将会非常舒服
	common.ConfBox = packr.New(configBoxName, configPath)
	// 获取实例(可创建多实例读取多个配置文件, 这里不做演示)
	v := viper.New()

	// 读取开发环境配置作为默认配置项
	readConfig(v, devConfig)
	// 将default中的配置全部以默认配置写入
	settings := v.AllSettings()
	for index, setting := range settings {
		v.SetDefault(index, setting)
	}
	// 读取当前go运行环境变量
	env := os.Getenv("GO_ENV")
	configName := ""
	if env == "test" {
		configName = testConfig
	} else if env == "prod" {
		configName = prodConfig
	}
	if configName != "" {
		// 读取不同环境中的差异部分
		readConfig(v, configName)
	}
	// 转换为结构体
	if err := v.Unmarshal(&common.Conf); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v", err))
	}

	fmt.Println("初始化配置文件完成")
}

func readConfig(v *viper.Viper, configFile string) {
	v.SetConfigType(configType)
	config, err := common.ConfBox.Find(configFile)
	if err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v", err))
	}
	// 加载配置
	if err = v.ReadConfig(bytes.NewReader(config)); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v", err))
	}
}
