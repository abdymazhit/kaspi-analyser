package scheduler

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"kaspi-analyser/internal/mongodb"
	"kaspi-analyser/pkg/httpClient"
	"kaspi-analyser/pkg/md5"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StartShopReviewsScheduler(ctx context.Context, dhm *mongodb.DBHandlerMongo, service *httpClient.Service) {
	startTime := time.Now()
	finish := make(chan struct{})

	go func() {
		log.Println("shop reviews scheduler: start shop reviews scheduler at", startTime)

		merchants, err := dhm.GetMerchantsFromOffers(ctx)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("shop reviews scheduler: error while getting merchants from offers:", err)
			finish <- struct{}{}
			return
		}

		for _, merchant := range merchants {
			page := 0
			for {
				response, err := service.SendJSONRequest(ctx, http.MethodGet, "https://kaspi.kz/shop/rest/misc/merchant/"+merchant.ID+"/reviews?limit=100&page="+strconv.Itoa(page), nil)
				if err != nil {
					log.Println("shop reviews scheduler: error while sending request for page", page, err)
					page++
					continue
				}

				data, ok := response["data"].([]interface{})
				if !ok {
					log.Println("shop reviews scheduler: error while parsing reviews response")
					page++
					continue
				}

				if len(data) == 0 {
					break
				}

				for _, review := range data {
					reviewMap, ok := review.(map[string]interface{})
					if !ok {
						log.Println("shop reviews scheduler: error while parsing review")
						continue
					}

					author, ok := reviewMap["author"].(string)
					if !ok {
						log.Println("shop reviews scheduler: error while parsing review author")
						continue
					}

					date, ok := reviewMap["date"].(string)
					if !ok {
						log.Println("shop reviews scheduler: error while parsing review date")
						continue
					}

					rating, ok := reviewMap["rating"].(float64)
					if !ok {
						log.Println("shop reviews scheduler: error while parsing review rating")
						continue
					}

					reviewMap["merchant_id"] = merchant.ID
					delete(reviewMap, "id")

					// create hash string by author, date, rating
					hashId := md5.GetMD5Hash(author + date + strconv.FormatFloat(rating, 'f', 0, 64))

					// save review to db
					if err = dhm.SaveShopReview(ctx, hashId, reviewMap); err != nil && err != mongo.ErrNoDocuments {
						log.Println("shop reviews scheduler: error while saving shop review:", err)
						continue
					}
				}
				page++
			}
		}

		finish <- struct{}{}
	}()

	for {
		select {
		case <-finish:
			log.Println("shop reviews scheduler: finish shop reviews scheduler in", time.Since(startTime))

			// sleep 3 min before start new shop reviews scheduler
			time.Sleep(3 * time.Minute)
			StartShopsScheduler(ctx, dhm)

			return
		}
	}
}
