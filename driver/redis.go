package driver

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {

	fmt.Println("Waiting for redis startup ...")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASS")

	if redisHost == "" {
		redisHost = "redis" // fallback for local development
	}
	if redisPort == "" {
		redisPort = "6379" // fallback for local development
	}

	var rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Username: "default", // remove this field when running locally
		Password: redisPass, // no password by default
		DB:       0,         // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("Redis connection failed:", err)
		return nil, err
	}

	log.Println("Connected to Redis!")
	return rdb, nil

}
