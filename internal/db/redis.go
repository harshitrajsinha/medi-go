package driver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis(host string, pass string, port string) (*redis.Client, error) {

	fmt.Println("Waiting for redis startup ...")

	var rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: "default", // remove this field when running locally
		Password: pass,      // no password by default
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
