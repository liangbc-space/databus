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
}

/*redis连接配置*/
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Prefix   string `yaml:"prefix"`
}

type Configuration struct {
	Port            uint           `yaml:"port"`
	DefaultPageSize uint           `yaml:"default_page_size"`
	JWTSecret       string         `yaml:"jwt_secret"`
	DbConfig        DatabaseConfig `yaml:"database"`
	RedisConfig     RedisConfig    `yaml:"redis"`
}

const (
	defaultPageSize = 10
	httpServerPort  = 8080
)

var SystemConfig *Configuration

func LoadConfiguration(path string) error {
	if SystemConfig != nil {
		return nil
	}

	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configData, &SystemConfig)
	if err != nil {
		return err
	}

	if SystemConfig.DefaultPageSize <= 0 {
		SystemConfig.DefaultPageSize = defaultPageSize
	}

	if SystemConfig.Port <= 0 {
		SystemConfig.Port = httpServerPort
	}

	return err

}
