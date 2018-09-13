package utils

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var _redis *redis.Client

func init() {

}

func Redis() *redis.Client {
	if _redis == nil {
		host := GetConfig().GetString("redis.host")
		port := GetConfig().GetString("redis.port")
		db := GetConfig().GetInt("redis.db")

		addr := fmt.Sprintf("%v:%v", host, port)
		options := redis.Options{
			Addr: addr,
			DB:   db,
		}

		if password := GetConfig().GetString("redis.password"); password != "" {
			log.Println("redis use auth", password)
			options.Password = password
		}

		client := redis.NewClient(&options)

		if _, err := client.Ping().Result(); err != nil {
			fmt.Println("ping error", addr, err)
		} else {
			fmt.Println("pong")
		}

		_redis = client
	}
	return _redis

}
