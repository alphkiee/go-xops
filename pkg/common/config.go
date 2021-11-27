package common

// 系统配置, 配置字段可参见yml注释
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便
type Configuration struct {
	System               SystemConfiguration `mapstructure:"system" json:"system"`
	Mysql                MysqlConfiguration  `mapstructure:"mysql" json:"mysql"`
	Jwt                  JwtConfiguration    `mapstructure:"jwt" json:"jwt"`
	Upload               UploadConfiguration `mapstructure:"upload" json:"upload"`
	Casbin               CasbinConfiguration `mapstructure:"casbin" json:"casbin"`
	PrometheusApiAddress PrometheusAddress   `mapstructure:"prometheusapiaddress" json:"prometheusapiaddress"`
	KubeConf             KubeConfiguration   `mapstructure:"kubeConf" json:"kubeConf"`
	DockerApi            DockerApiUrl        `mapstructure:"dockerapi" json:"docker_api"`
	HarborApi            HarborApiUrl        `mapstructure:"harborapi" json:"harbor_api"`
	ZapLogPath           ZapLog              `mapstructure:"zapLog" json:"zap_log"`
}

type SystemConfiguration struct {
	UrlPathPrefix   string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	AppMode         string `mapstructure:"app-mode" json:"appMode"`
	Transaction     bool   `mapstructure:"transaction" json:"transaction"`
	Port            int    `mapstructure:"port" json:"port"`
	OperationLogKey string `mapstructure:"operation-log-key" json:"operationLogKey"`
}

type MysqlConfiguration struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}

type CasbinConfiguration struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath"`
	LoadDelay int    `mapstructure:"load-delay" json:"loadDelay"`
}

type JwtConfiguration struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type UploadConfiguration struct {
	SaveDir       string `mapstructure:"save-dir" json:"saveDir"`
	SingleMaxSize uint   `mapstructure:"single-max-size" json:"singleMaxSize"`
}

type PrometheusAddress struct {
	Address string `mapstructure:"address" json:"address"`
}

type KubeConfiguration struct {
	Path string `mapstructure:"path" json:"path"`
}

// Docker api url 地址
type DockerApiUrl struct {
	Url string `mapstructure:"url" json:"url"`
}

// Harbor api url 地址
type HarborApiUrl struct {
	Url      string `mapstructure:"url" json:"url"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

// zap log
type ZapLog struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}
