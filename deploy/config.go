package deploy

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
	"microsvc/consts"
)

// XConfig 是主配置结构体
type XConfig struct {
	Svc   string             `mapstructure:"svc"`
	Env   consts.Environment `mapstructure:"env"`
	Mysql map[string]*Mysql  `mapstructure:"mysql"`
	Redis map[string]*Redis  `mapstructure:"redis"`
}

var XConf = &XConfig{}

func init() {
	XConf.Env = readEnv()

	// 设置配置文件名（不包含扩展名）
	viper.SetConfigName("config")
	// 设置配置文件所在的路径（可选，默认为当前目录）
	viper.AddConfigPath("deploy/" + string(XConf.Env))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config file: %s\n", err))
		return
	}
	if err := viper.Unmarshal(XConf); err != nil {
		panic(fmt.Sprintf("Error Unmarshal config: %s\n", err))
		return
	}

	_, _ = pp.Printf("********* init Config OK *********\n%+v\n", XConf)
}

type DBname string

type Mysql struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	GormArgs string `mapstructure:"gorm_args"`
}

type Redis struct {
	DB       int8   `mapstructure:"db"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
}
