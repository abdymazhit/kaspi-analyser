package scheduler

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"kaspi-analyser/internal/mongodb"
	"log"
	"time"
)

func StartShopsScheduler(ctx context.Context, dhm *mongodb.DBHandlerMongo) {
	startTime := time.Now()
	finish := make(chan struct{})

	go func() {
		log.Println("scheduler: start shop scheduler at", startTime)

		merchants, err := dhm.GetMerchantsFromOffers(ctx)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("scheduler: error while getting merchants from offers:", err)
			finish <- struct{}{}
			return
		}

		// merchants to []interface
		shops := make([]interface{}, len(merchants))
		for i, merchant := range merchants {
			shops[i] = merchant
		}

		// save shops to db
		if err = dhm.SaveShops(ctx, shops); err != nil && err != mongo.ErrNoDocuments && !mongo.IsDuplicateKeyError(err) {
			log.Println("scheduler: error while saving shop:", err)
			finish <- struct{}{}
			return
		}

		finish <- struct{}{}
	}()

	for {
		select {
		case <-finish:
			log.Println("scheduler: finish shop scheduler in", time.Since(startTime))

			// sleep 3 min before start new shops scheduler
			time.Sleep(3 * time.Minute)
			StartShopsScheduler(ctx, dhm)

			return
		}
	}
}
