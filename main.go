package main

import (
	"flag"
	"fmt"
	"github.com/liangbc-space/databus/models"
	"github.com/liangbc-space/databus/system"
	"github.com/liangbc-space/databus/utils"
	"github.com/liangbc-space/databus/utils/exception"
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
	models.InitDB()

	//	初始化redis连接池
	//utils.InitRedis()

	//	初始化elasticsearch链接
	utils.InitElasticsearch()
}

func main() {

	defer models.DB.Close()

	//mysql_elasticsearch.Run()
	exception.Throw("error2", 0,nil)
	return
	exception.Try(func() {
		//  指定了异常代码为2，错误信息为error2
		exception.Throw("error2", 0,exception.Exception{})
	}).Catch(exception.Exception{}, func(e interface{}) {
		fmt.Println(e)
	}).Finally(func() {
		fmt.Println(123123)
	})

}
