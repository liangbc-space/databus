package main

import (
	"demo/models"
	"demo/routers"
	"demo/system"
	"demo/utils"
	"flag"
	"fmt"
)

func init() {
	//	初始化配置
	configPath := flag.String("systemConfig", "conf/conf.yaml", "system config file path")
	flag.Parse()

	err := system.LoadConfiguration(*configPath)
	if err != nil {
		panic(err)
		return
	}

	//	初始化数据库
	_, err = models.InitDB()
	if err != nil {
		panic(err)
		return
	}

	//	初始化redis连接池
	utils.InitRedis()
}

func main() {
	//	释放数据库连接
	defer models.DB.Close()

	request := routers.InitRouter()

	if err := request.Run(fmt.Sprintf(":%d", system.SystemConfig.Port)); err != nil {
		panic(err)
	}

}
