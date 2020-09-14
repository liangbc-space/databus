package main

import (
	"databus/models"
	mysql_elasticsearch "databus/mysql-elasticsearch"
	"fmt"
	"reflect"

	//"databus/routers"
	"databus/system"
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
	/*if _, err := models.InitDB(); err != nil {
		panic(err)
		return
	}

	//	初始化redis连接池
	utils.InitRedis()

	//	初始化elasticsearch链接
	utils.InitElasticsearch()*/
}

func main() {

	a := struct {
		Id	uint32
		Name	string
	}{56, "测试"}
	fmt.Println(a)

	b := make(map[string]interface{})

	if _,ok:=b["category_ids"];!ok {
		b["category_ids"] = make([]uint32,0)
	}

	rType := reflect.TypeOf(b["category_ids"])
	rValue := reflect.ValueOf(b["category_ids"])

	fmt.Println(rType.Kind())
	switch rType.Kind() {
	case reflect.Slice,reflect.Array:
		fmt.Println(rValue.Len())
		/*fmt.Println(reflect.ValueOf(a).FieldByName("Id"))
		s := make([]reflect.Value,0 )
		s = append(s,reflect.ValueOf(a).FieldByName("Id"))
		value := reflect.Append(rValue, s...)
		rValue.Elem().Set(value)*/
	}

	fmt.Println(b["category_ids"])

	return
	defer models.DB.Close()

	mysql_elasticsearch.Run()

}
