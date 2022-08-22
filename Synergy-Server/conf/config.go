package conf

import (
	"github.com/pelletier/go-toml"
	json "github.com/wulinlw/jsonAmbiguity"
	"io/ioutil"
	"log"
	"time"
)

type Redis struct {
	Port     int    `json:"port" toml:"port"`
	Host     string `json:"host" toml:"host"`
	Family   int    `json:"family" toml:"family"`
	Password string `json:"password" toml:"password"`
	Db       int    `json:"db" toml:"db"`
}

type Mysql struct {
	User     string `json:"user" toml:"user"`
	Host     string `json:"host" toml:"host"`
	Port     int    `json:"port" toml:"port"`
	Database string `json:"database" toml:"database"`
	Password string `json:"password" toml:"password"`
}

type Logger struct {
	//日志文件保存天数
	LogExpire int    `json:"logExpire" toml:"logExpire"`
	LogDir    string `json:"logDir" toml:"logDir"`
	LogName   string `json:"logName" toml:"logName"`
	LogLevel  int    `json:"logLevel" toml:"logLevel"`
}

type StorageConfig struct {
	//第三方云存储的 access key。建议提供只写权限的访问密钥（当vendor=100时，此处可为""）。
	AccessKey string `json:"accessKey"`
	//第三方云存储指定的地区信息
	Region int `json:"region,omitempty"`
	//第三方云存储的 bucket（当vendor=100时，bucket应为访问私有存储的接口地址）。
	Bucket string `json:"bucket"`
	//第三方云存储的 secret key（当vendor=100时，此处可为""）。
	SecretKey string `json:"secretKey"`
	//第三方云存储供应商。
	Vendor int `json:"vendor"`
	//由多个字符串组成的数组，指定录制文件在第三方云存储中的存储位置
	FileNamePrefix []string `json:"fileNamePrefix"`
}

type ServerConfig struct {
	//进程名字
	ProcessName string `json:"processName" toml:"processName"`
	//服务监听端口
	Host string `json:"host" toml:"host"`
	//端口
	Port int `json:"port" toml:"port"`
	//进程保存目录
	PidDir string `json:"pidDir" toml:"pidDir"`

	//token有效期（单位：秒）
	TokenExpire time.Duration `json:"tokenExpire" toml:"tokenExpire"`
	//rtc token有效期（单位：秒）
	RtcTokenExpire time.Duration `json:"rtcTokenExpire" toml:"rtcTokenExpire"`
	//rtm token有效期（单位：秒）
	RtmTokenExpire time.Duration `json:"rtmTokenExpire" toml:"rtmTokenExpire"`
	//查询用户心跳包间隔时间（单位：秒）
	OnlineInterval time.Duration `json:"onlineInterval" toml:"onlineInterval"`

	//web资源目录
	WebSource string `json:"webSource" toml:"webSource"`

	//redis配置文件
	Redis Redis `json:"redis" toml:"redis"`
	//mysql配置文件
	Mysql Mysql `json:"mysql" toml:"mysql"`
	//第三方云存储的配置文件
	StorageConfig StorageConfig `json:"storageConfig" toml:"storageConfig"`

	//logger配置文件
	Logger     Logger `json:"logger" toml:"logger"`
	PkgSrcDir  string `json:"pkgSrcDir" toml:"pkgSrcDir"`
	RtmAddr    string `json:"rtmAddr" toml:"rtmAddr"`
	RtmWsPort  int    `json:"rtmWsPort" toml:"rtmWsPort"`
	RtmCliPort int    `json:"rtmCliPort" toml:"rtmCliPort"`

	//开发者通知服务用于签名的secret
	DevpSecret string `json:"devpSecret" toml:"devpSecret"`
	//用户ID
	CustomerId string `json:"customerId" toml:"customerId"`
	//用户证书(certificate)
	CustomerCert string `json:"customerCert" toml:"customerCert"`

	//oss存储路径
	OssDir string `json:"ossDir" toml:"ossDir"`
	//oss存储系统前缀
	HttpOssPrefix string `json:"httpOssPrefix" toml:"httpOssPrefix"`

	//录像服务接口前缀
	HttpVodPrefix string `json:"httpVodPrefix" toml:"httpVodPrefix"`
	//录像文件播放下载地址前缀
	HttpVodFilePrefix string `json:"httpVodFilePrefix" toml:"httpVodFilePrefix"`
	//录像文件保存天数
	VodValidDays int64 `json:"vodValidDays" toml:"vodValidDays"`
	//Bucket对应的Endpoint
	OssEndpoint string `json:"ossEndpoint" toml:"ossEndpoint"`

	//appId
	AppId string `json:"appId" toml:"appId"`
	//appToken
	AppToken string `json:"appToken" toml:"appToken"`
}

type AllConfig struct {
	UseConf string       `json:"useConf" toml:"useConf"`
	DConfig ServerConfig `json:"dConfig" toml:"dConfig"`
	RConfig ServerConfig `json:"rConfig" toml:"rConfig"`
	PConfig ServerConfig `json:"pConfig" toml:"pConfig"`
	TConfig ServerConfig `json:"tConfig" toml:"tConfig"`
}

func NewAllConfig() *AllConfig {
	return &AllConfig{}
}

func (ac *AllConfig) LoadJsonConfig(confPath string) (ok bool) {
	if confPath == "" {
		log.Fatal("Params not allowed, confPath is null")
		return false
	}
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal("ReadFile failed, err: ", err)
		return false
	}
	err = json.Unmarshal(data, ac)
	if err != nil {
		log.Fatal("Unmarshal failed, err: ", err)
		return false
	}
	return true
}

func (ac *AllConfig) LoadTomlConfig(confPath string) (ok bool) {
	if confPath == "" {
		log.Fatal("Params not allowed, confPath is null")
		return false
	}
	conf, err := toml.LoadFile(confPath)
	if err != nil {
		log.Fatal("ReadFile failed, err: ", err)
		return false
	}
	err = conf.Unmarshal(ac)
	if err != nil {
		log.Fatal("Unmarshal failed, err: ", err)
		return false
	}
	return true
}
