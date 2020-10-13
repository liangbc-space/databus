package system

import (
	"github.com/liangbc-space/databus/utils/exception"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type LoggerConfig struct {
	Level        string `yaml:"level"`
	LogPath      string `yaml:"log_path"`
	LogValidDays uint   `yaml:"log_valid_days"`
}

/*mysql连接配置*/
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Prefix   string `yaml:"prefix"`
	Debug    bool   `yaml:"debug"`
	Pool     struct {
		MaxIdle uint `yaml:"max_idle"`
		MaxOpen uint `yaml:"max_open"`
	} `yaml:"pool"`
}

/*redis连接配置*/
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Prefix   string `yaml:"prefix"`
}

/*kafka连接配置*/
type KafkaConfig struct {
	Brokers             []string `yaml:"brokers"`
	BrokerAddressFamily string   `yaml:"broker_address_family"`
	ConsumerLogs        bool     `yaml:"consumer_logs"`
}

/*elasticsearch连接配置*/
type ElasticsearchConfig struct {
	Urls     []string `yaml:"urls"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type configuration struct {
	AppName             string              `yaml:"app_name"`
	Port                uint                `yaml:"port"`
	DefaultPageSize     uint                `yaml:"default_page_size"`
	JWTSecret           string              `yaml:"jwt_secret"`
	Debug               bool                `yaml:"debug"`
	LoggerConfig        LoggerConfig        `yaml:"logger"`
	DbConfig            DatabaseConfig      `yaml:"mysql"`
	RedisConfig         RedisConfig         `yaml:"redis"`
	KafkaConfig         KafkaConfig         `yaml:"kafka"`
	ElasticsearchConfig ElasticsearchConfig `yaml:"elasticsearch"`
}

const (
	defaultPageSize = 10
	httpServerPort  = 8080
)

var ApplicationCfg *configuration

func LoadConfiguration(path string) {
	if ApplicationCfg != nil {
		return
	}

	configData, err := ioutil.ReadFile(path)
	if err != nil {
		exception.Throw("初始化系统配置失败："+err.Error(), 1)
	}

	err = yaml.Unmarshal(configData, &ApplicationCfg)
	if err != nil {
		exception.Throw("初始化系统配置失败："+err.Error(), 1)
	}

	if ApplicationCfg.DefaultPageSize <= 0 {
		ApplicationCfg.DefaultPageSize = defaultPageSize
	}

	if ApplicationCfg.Port <= 0 {
		ApplicationCfg.Port = httpServerPort
	}

	if ApplicationCfg.LoggerConfig == (LoggerConfig{}) {
		if ApplicationCfg.Debug {
			ApplicationCfg.LoggerConfig.Level = "debug"
		} else {
			ApplicationCfg.LoggerConfig.Level = "warn"
		}

		ApplicationCfg.LoggerConfig.LogPath = "logs/debug.log"
		ApplicationCfg.LoggerConfig.LogValidDays = 10
	}

}
