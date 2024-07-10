package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx context.Context = context.Background()

func ConnectCache() {
	serviceURI := os.Getenv("REDIS")

	addr, err := redis.ParseURL(serviceURI)

	if err != nil {
		log.Fatalf(err.Error())
	}

	rdb := redis.NewClient(addr)

	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		_, err := rdb.Ping(ctx).Result()

		if err == nil {
			fmt.Println("Redis Cached connected on PORT :: 5500")
			break
		}

		if i < maxRetries-1 {
			log.Printf("Failed to connect to redis, retrying ... (%d/%d)", i+1, maxRetries)
			backOff := time.Duration(math.Pow(float64(i), 2)) * time.Second
			time.Sleep(backOff)
		} else {
			log.Fatalf("Failed to connect to redis, after %d attempts: %v", maxRetries, err)
		}
	}

	RedisClient = rdb
}

func UpdateHash(ctx context.Context,cacheKey string,fieldKey string,value interface{}) error {
    jsonMarshal,err := json.Marshal(value);

	if err != nil {
		return err
	}

	cacheSet := RedisClient.HSet(ctx,cacheKey,fieldKey,string(jsonMarshal))

	if cacheSet.Err() != nil {
		return cacheSet.Err()
	}

	return nil
}

func GetHashKeys(ctx context.Context,cacheKey string,fieldKey string,data interface{}) ([]string,error) {
	cacheGet := RedisClient.HKeys(ctx,cacheKey);

	if cacheGet.Err() != nil {
		return nil,cacheGet.Err()
	}

	return  cacheGet.Val(),nil
}
