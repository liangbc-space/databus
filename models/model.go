package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/liangbc-space/databus/system"
	"time"
)

type BaseModel struct {
	Id         uint  `gorm:"primary_key" json:"id" redis:"id"`
	CreateTime int64 `json:"create_time" redis:"create_time"`
	UpdateTime int64 `json:"update_time" redis:"update_time"`
}

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	var err error

	DbConfig := system.ApplicationCfg.DbConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", DbConfig.Username, DbConfig.Password, DbConfig.Host, DbConfig.Port, DbConfig.Dbname)

	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	//	全局禁用负数表名
	DB.SingularTable(true)
	//	默认表名规则
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return DbConfig.Prefix + defaultTableName
	}
	DB.LogMode(DbConfig.Debug)

	DB.DB().SetMaxIdleConns(int(DbConfig.Pool.MaxIdle)) //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭
	DB.DB().SetMaxOpenConns(int(DbConfig.Pool.MaxOpen)) //设置数据库连接池最大连接数

	return DB, err

}

func (model *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("create_time", time.Now().Unix()); err != nil {
		return err
	}

	if err := scope.SetColumn("update_time", time.Now().Unix()); err != nil {
		return err
	}
	return nil
}

func (model *BaseModel) BeforeUpdate(scope *gorm.Scope) error {
	if err := scope.SetColumn("update_time", time.Now().Unix()); err != nil {
		return err
	}
	return nil
}
