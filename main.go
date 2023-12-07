package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"poc-growthbook/internal/featureflag"
	"poc-growthbook/pkg/handler"
	"poc-growthbook/pkg/middleware"
)

func main() {
	gbs := featureflag.NewGrowthBookService(os.Getenv("GROWTHBOOK_URL"))

	http.HandleFunc("/callback", handler.Callback(gbs))
	http.Handle("/", middleware.InjectUserData(handler.Home(gbs)))

	if err := http.ListenAndServe(":8080", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}
