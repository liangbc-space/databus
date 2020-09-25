package utils

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/liangbc-space/databus/system"
	"reflect"
	"time"
)

const (
	MAX_IDLE             = 200  //	最大空闲连接数
	MAX_ACTIVE           = 500  //	最大连接数
	IDLE_TIMEOUT_SEC     = 240  //	最大空闲连接时间，秒
	DIAL_READ_TIMEOUT    = 1000 //	最大读数据超时时间	毫秒
	DIAL_WRITE_TIMEOUT          //	最大写数据超时时间	毫秒
	DIAL_CONNECT_TIMEOUT        //	最大连接数据超时时间	毫秒
)

type redisClient struct {
	Pool *redis.Pool
}

var (
	RedisClient *redisClient
	redisConfig system.RedisConfig
)

func InitRedis() {
	redisConfig = system.ApplicationCfg.RedisConfig

	RedisClient = new(redisClient)
	RedisClient.Pool = &redis.Pool{
		MaxIdle:     MAX_IDLE,
		MaxActive:   MAX_ACTIVE,
		IdleTimeout: time.Duration(IDLE_TIMEOUT_SEC),
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
				redis.DialReadTimeout(time.Duration(DIAL_READ_TIMEOUT)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(DIAL_WRITE_TIMEOUT)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(DIAL_CONNECT_TIMEOUT)*time.Millisecond),
				redis.DialDatabase(redisConfig.Db),
				redis.DialPassword(redisConfig.Password),
			)
		},
	}
}

func (redisClient *redisClient) Exec(cmd string, key string, args ...interface{}) (reply interface{}, err error) {

	conn := redisClient.Pool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	defer conn.Close()

	key = redisConfig.Prefix + ":" + key
	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, v := range args {
			switch reflect.ValueOf(v).Kind() {
			case reflect.Map, reflect.Struct, reflect.Ptr, reflect.Slice:
				params = append(params, redis.Args{}.AddFlat(v)...)
			default:
				params = append(params, v)
			}

		}
	}

	return conn.Do(cmd, params...)
}
