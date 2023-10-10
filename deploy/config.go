package deploy

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
	"microsvc/enums"
	"microsvc/enums/svc"
	"microsvc/util"
	"os"
	"path/filepath"
)

// XConfig 是主配置结构体
type XConfig struct {
	Svc                   svc.Svc           `mapstructure:"svc"` // set by this.svcConf
	Env                   enums.Environment `mapstructure:"env"`
	Mysql                 map[string]*Mysql `mapstructure:"mysql"`
	Redis                 map[string]*Redis `mapstructure:"redis"`
	SimpleSdHttpPort      int               `mapstructure:"simplesd_http_port"`       // 本地简单注册中心的固定端口，可在配置修改
	SvcTokenSignKey       string            `mapstructure:"svc_token_sign_key"`       // 微服务鉴权token使用的key
	AdminTokenSignKey     string            `mapstructure:"admin_token_sign_key"`     // Admin鉴权token使用的key
	SensitiveInfoCryptKey string            `mapstructure:"sensitive_info_crypt_key"` // 敏感信息加密key（如手机号、身份证等）

	// 私有字段
	gRPCPort int
	httpPort int

	// 接管svc的配置
	svcConf SvcConfImpl
}

func (x *XConfig) GetSvc() string {
	return x.Svc.Name()
}

func (x *XConfig) SetGRPC(port int) {
	x.gRPCPort = port
}

func (x *XConfig) SetHTTP(port int) {
	x.httpPort = port
}

func (s *XConfig) RegGRPCBase() (name string, addr string, port int) {
	return s.Svc.Name(), "", s.gRPCPort
}

func (s *XConfig) RegGRPCMeta() map[string]string {
	return nil
}

func (s *XConfig) GetSvcConf() SvcConfImpl {
	return s.svcConf
}

func (s *XConfig) GetConfDir(subPath ...string) string {
	return filepath.Join(append([]string{"deploy", s.Env.S()}, subPath...)...)
}

type Initializer func(cc *XConfig)

var XConf = &XConfig{}

var _ SvcListenPortSetter = new(XConfig)
var _ RegisterSvc = new(XConfig)

func Init(svc svc.Svc, svcConfVar SvcConfImpl) {
	XConf.Svc = svc
	XConf.Env = readEnv()

	// 设置配置文件名（不包含扩展名）
	viper.SetConfigName("config")
	// 设置配置文件所在的路径（可选，默认为当前目录）
	viper.AddConfigPath(XConf.GetConfDir())
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
	for dbname, v := range XConf.Redis {
		v.DBname = DBname(dbname)
	}
	_, _ = pp.Printf("\n************* init Share-Config OK *************\n%+v\n", XConf)

	// ------------- 下面读取svc专有配置 -------------------

	if svcConfVar != nil {
		//wd, _ := os.Getwd()
		//println("getwd", wd)
		svcConfFile, err := os.Open(fmt.Sprintf("service/%s/deploy/%s/config.yaml", svc, XConf.Env))
		util.AssertNilErr(err)

		err = viper.ReadConfig(svcConfFile)
		util.AssertNilErr(err)

		err = viper.Unmarshal(svcConfVar)
		util.AssertNilErr(err)

		logLv := readLogLevelFromEnvVar()
		if logLv != "" {
			svcConfVar.OverrideLogLevel(logLv)
			_, _ = pp.Printf("************* read log level from env: %s\n", logLv)
		}

		if svc != svcConfVar.GetSvc() {
			panic(fmt.Sprintf("%s not match svc name:%s in config file", svc, svcConfVar.GetSvc()))
		}
		_, _ = pp.Printf("\n************* init Svc-Config OK *************\n%+v\n", svcConfVar)

		// svc conf 嵌入主配置
		XConf.svcConf = svcConfVar
	}
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
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s&timeout=3s", m.User, m.Password, m.Host, m.Port, m.DBname, m.GormArgs)
}

type Redis struct {
	DBname
	DB       int    `mapstructure:"db"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
}
