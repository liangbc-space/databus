package system

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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
}

/*elasticsearch连接配置*/
type ElasticsearchConfig struct {
	Urls     []string `yaml:"urls"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type Configuration struct {
	Port                uint                `yaml:"port"`
	DefaultPageSize     uint                `yaml:"default_page_size"`
	JWTSecret           string              `yaml:"jwt_secret"`
	DbConfig            DatabaseConfig      `yaml:"mysql"`
	RedisConfig         RedisConfig         `yaml:"redis"`
	KafkaConfig         KafkaConfig         `yaml:"kafka"`
	ElasticsearchConfig ElasticsearchConfig `yaml:"elasticsearch"`
}

const (
	defaultPageSize = 10
	httpServerPort  = 8080
)

var ApplicationCfg *Configuration

func LoadConfiguration(path string) error {
	if ApplicationCfg != nil {
		return nil
	}

	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configData, &ApplicationCfg)
	if err != nil {
		return err
	}

	if ApplicationCfg.DefaultPageSize <= 0 {
		ApplicationCfg.DefaultPageSize = defaultPageSize
	}

	if ApplicationCfg.Port <= 0 {
		ApplicationCfg.Port = httpServerPort
	}

	return err

}
