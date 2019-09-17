package cache

import (
	"CDcoding2333/scaffold/conf"
	"CDcoding2333/scaffold/utils/nosql"
)

var globalRedisStore *nosql.RedisStore

//GetRedis get redis store
func GetRedis() *nosql.RedisStore {
	return globalRedisStore
}

//InitRedis ...
func InitRedis(config *conf.RedisConfig) (err error) {
	globalRedisStore, err = nosql.NewRedisStore(config.Host, config.Password, config.Port, config.DB)
	return
}
