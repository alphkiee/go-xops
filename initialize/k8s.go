package initialize

import (
	"go-xops/pkg/common"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func InitKubeConf() {
	c := common.Conf.KubeConf
	// 加载k8s配置文件，生成Config对象
	config, err := clientcmd.BuildConfigFromFlags("", c.Path)
	if err != nil {
		panic(err)
	}
	common.ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}
