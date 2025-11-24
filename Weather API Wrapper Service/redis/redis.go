package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RDB *redis.Client

func ConnectToRedis() {
	Address := getEnv("REDIS_HOST", "localhost:6379")
	Password := getEnv("REDIS_PASSWORD", "")
	Port := getEnv("REDIS_PORT", "6379")

	RDB = redis.NewClient(&redis.Options{
		Addr:     Address + ":" + Port,
		Password: Password,
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to redis ")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
