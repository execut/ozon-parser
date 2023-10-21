package main

import "github.com/redis/go-redis/v9"

var rdb *redis.Client

func SetCachedValue(url string, dataJson string) {
	err := rdb.Set(ctx, url, dataJson, 0).Err()
	if err != nil {
		panic(err)
	}
}

func GetCachedValue(url string) (string, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	val, err := rdb.Get(ctx, url).Result()
	if err != redis.Nil && err != nil {
		panic(err)
	}

	return val, err
}
