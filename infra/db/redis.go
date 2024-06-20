package db

import (
	"context"
	"fmt"
	"log"

	"fashion-api/infra/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {

	appConfig := config.NewAppConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", appConfig.RedisHost, appConfig.RedisPort),
		Password: appConfig.RedisPass,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return rdb
}
