package main

import (
    "context"
    "github.com/redis/go-redis/v9"
    "os"
    "time"
)

var rdb *redis.Client = nil
var redisCtx context.Context

func getSecondOfDayBetween() time.Duration {
    t := time.Now()
    return time.Duration(float64(60*60*24-60*60*t.Hour()-60*t.Minute()-t.Second()) * float64(time.Second))
}

func SetCachedValue(url string, dataJson string) {
    initRedisClient()
    err := rdb.Set(redisCtx, url, dataJson, getSecondOfDayBetween()).Err()
    if err != nil {
        panic(err)
    }
}

func GetCachedValue(url string) (string, error) {
    initRedisClient()
    //return "", redis.Nil

    val, err := rdb.Get(redisCtx, url).Result()
    if err != redis.Nil && err != nil {
        panic(err)
    }

    return val, err
}

func initRedisClient() {
    if rdb == nil {
        redisCtx = context.Background()
        rdb = redis.NewClient(&redis.Options{
            Addr:     os.Getenv("REDIS_HOST") + ":6379",
            Password: os.Getenv("REDIS_PASSWORD"),
            DB:       0,
        })
    }
}
