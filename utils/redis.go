package utils

import (
	"fmt"

	"github.com/go-redis/redis"
)

var Redis *redis.Client

func init() {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", GetConfig().Redis.Host, GetConfig().Redis.Port),
		Password: GetConfig().Redis.Password, // no password set
		DB:       0,                          // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		fmt.Println("ping error", err)
	}

	Redis = client
}
