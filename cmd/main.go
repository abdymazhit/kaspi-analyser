package main

import (
	"context"
	"kaspi-analyser/internal/mongodb"
	"kaspi-analyser/internal/scheduler"
	"kaspi-analyser/pkg/httpClient"
	"log"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init mongodb
	dhm, err := mongodb.NewDBHandlerMongo(ctx, mongodb.DBConfigMongo{
		URI: os.Getenv("KA_MONGO_URI"),
	})
	if err != nil {
		log.Fatal("error while connecting to mongodb:", err)
		return
	}

	service := httpClient.NewService()

	// start products scheduler
	go func() {
		scheduler.StartProductsScheduler(ctx, dhm, service)
	}()

	// start shops scheduler
	go func() {
		scheduler.StartShopsScheduler(ctx, dhm)
	}()

	// start shop reviews scheduler
	go func() {
		scheduler.StartShopReviewsScheduler(ctx, dhm, service)
	}()

	<-ctx.Done()
	log.Println("main: context is done")
}
