package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func CreateClient(ctx context.Context, dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"), // no password set
		DB:       dbNo,                 // use default DB
	})

	// Test the connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return rdb
}
