package utils

import (
	"github.com/liangbc-space/databus/system"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"time"
)

const (
	maxIdle            = 200  //	最大空闲连接数
	maxActive          = 500  //	最大连接数
	idleTimeoutSec     = 240  //	最大空闲连接时间，秒
	dialReadTimeout    = 1000 //	最大读数据超时时间	毫秒
	dialWriteTimeout          //	最大写数据超时时间	毫秒
	dialConnectTimeout        //	最大连接数据超时时间	毫秒
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
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeoutSec),
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
				redis.DialReadTimeout(time.Duration(dialReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(dialWriteTimeout)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(dialConnectTimeout)*time.Millisecond),
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
