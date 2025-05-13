package driver

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis(hostname string) (*redis.Client, error) {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     hostname,
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
