package utils

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var Redis *redis.Client

func init() {

	host := GetConfig().Redis.Host
	port := GetConfig().Redis.Port

	if v := os.Getenv("REDIS_HOST"); v != "" {
		host = v
	}

	if v := os.Getenv("REDIS_PORT"); v != "" {
		port = v
	}

	addr := fmt.Sprintf("%v:%v", host, port)
	options := redis.Options{
		Addr: addr,
		DB:   0, // use default DB
	}

	if GetConfig().Redis.Password != "" {
		options.Password = GetConfig().Redis.Password
	}

	client := redis.NewClient(&options)

	if _, err := client.Ping().Result(); err != nil {
		fmt.Println("ping error", addr, err)
	} else {
		fmt.Println("pong")
	}

	Redis = client
}
