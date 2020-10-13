package main

import (
	"flag"
	"fmt"
	"github.com/liangbc-space/databus/models"
	mysql_elasticsearch "github.com/liangbc-space/databus/mysql-elasticsearch"
	"github.com/liangbc-space/databus/system"
	"github.com/liangbc-space/databus/utils"
	"github.com/liangbc-space/databus/utils/exception"
	"os"
)

func init() {
	exception.Try(func() {
		//	初始化配置
		configPath := flag.String("systemConfig", "conf/application.yaml", "system config file path")
		flag.Parse()
		system.LoadConfiguration(*configPath)

		//	初始化mysql数据库
		models.InitDB()

		//	初始化redis连接池
		//utils.InitRedis()

		//	初始化elasticsearch链接
		utils.InitElasticsearch()
	}).Catch(func(ex exception.Exception) {
		fmt.Printf("程序执行异常：%s	文件：%s:%d\n",
			ex.Message(), ex.File(), ex.Line(),
		)
		os.Exit(ex.Code())
	})

}

func main() {

	defer models.DB.Close()

	exception.Try(func() {
		mysql_elasticsearch.Run()
	}).Catch(func(ex exception.Exception) {
		fmt.Println(ex.Message())
	})

}
