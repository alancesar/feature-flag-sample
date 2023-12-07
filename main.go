package main

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"poc-growthbook/internal/cache"
	"poc-growthbook/internal/featureflag"
	"poc-growthbook/pkg/handler"
	"poc-growthbook/pkg/middleware"
	"strconv"
)

func main() {
	redisDBEnv := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBEnv)
	if err != nil {
		log.Fatalln("while reading redis database:", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
		DB:   redisDB,
	})

	gbs := featureflag.NewGrowthBookService(os.Getenv("GROWTHBOOK_URL"), cache.NewRedis(redisClient))

	http.HandleFunc("/callback", handler.Callback(gbs))
	http.Handle("/", middleware.InjectUserData(handler.Home(gbs)))

	if err := http.ListenAndServe(":8080", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}
