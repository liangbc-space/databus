package main

import (
	"github.com/liangbc-space/databus/models"
	mysql_elasticsearch "github.com/liangbc-space/databus/mysql-elasticsearch"
	"github.com/liangbc-space/databus/system"
	"github.com/liangbc-space/databus/utils"
	"flag"
)

func init() {
	//	初始化配置
	configPath := flag.String("systemConfig", "conf/application.yaml", "system config file path")
	flag.Parse()

	if err := system.LoadConfiguration(*configPath); err != nil {
		panic(err)
		return
	}

	//	初始化mysql数据库
	if _, err := models.InitDB(); err != nil {
		panic(err)
		return
	}

	//	初始化redis连接池
	//utils.InitRedis()

	//	初始化elasticsearch链接
	utils.InitElasticsearch()
}

func main() {

	defer models.DB.Close()

	mysql_elasticsearch.Run()

}
