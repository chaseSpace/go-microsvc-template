package deploy

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
	"microsvc/enums"
	"microsvc/util"
	"os"
)

// XConfig 是主配置结构体
type XConfig struct {
	Svc   string            `mapstructure:"svc"`
	Env   enums.Environment `mapstructure:"env"`
	Mysql map[string]*Mysql `mapstructure:"mysql"`
	Redis map[string]*Redis `mapstructure:"redis"`

	// 接管svc的配置
	svcConf SvcConfImpl
}

func (x XConfig) GetSvcConf() SvcConfImpl {
	return x.svcConf
}

type Initializer func(cc *XConfig)

type SvcConfImpl interface {
	GetLogLevel() string
}

var XConf = &XConfig{}

func Init(svc string, svcConfVar SvcConfImpl) {

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
	for dbname, v := range XConf.Mysql {
		v.DBname = DBname(dbname)
	}
	_, _ = pp.Printf("\n************* init Share-Config OK *************\n%+v\n", XConf)

	// ------------- 下面读取svc专有配置 -------------------

	svcConfFile, err := os.Open(fmt.Sprintf("service/%s/deploy/%s/config.yaml", svc, XConf.Env))
	util.AssertNilErr(err)

	err = viper.ReadConfig(svcConfFile)
	util.AssertNilErr(err)

	err = viper.Unmarshal(&svcConfVar)
	util.AssertNilErr(err)
	util.AssertNotNil(svcConfVar)

	_, _ = pp.Printf("\n************* init Svc-Config OK *************\n%+v\n", svcConfVar)

	// svc conf 嵌入主配置
	XConf.svcConf = svcConfVar
}

type DBname string

type Mysql struct {
	DBname
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	GormArgs string `mapstructure:"gorm_args"`
}

func (m Mysql) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", m.User, m.Password, m.Host, m.Port, m.DBname, m.GormArgs)
}

type Redis struct {
	DB       int8   `mapstructure:"db"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
}
