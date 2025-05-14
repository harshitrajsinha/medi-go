package driver

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitRedis(hostname string) (*redis.Client, error) {

	_ = godotenv.Load()

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "redis" // fallback for local development
	}
	if redisPort == "" {
		redisPort = "6379" // fallback for local development
	}
	var rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // no password by default
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("Redis connection failed:", err)
		return nil, err
	}

	log.Println("Connected to Redis!")
	return rdb, nil

}
