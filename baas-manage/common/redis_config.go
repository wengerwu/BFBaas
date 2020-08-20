package common

import (
	"github.com/go-redis/redis"
	"github.com/paybf/baasmanager/baas-gateway/config"
)

var RedisClient *redis.Client

func ConnRedis() {
	url := MyDESDecrypt(config.Config.GetString("Redis.Url"))
	password:=config.Config.GetString("Redis.Password")
	if password!=""{
		password=MyDESDecrypt(password)
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0,
	})

	RedisClient.Ping().Result()
}
