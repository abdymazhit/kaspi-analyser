package scheduler

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"kaspi-analyser/internal/mongodb"
	"kaspi-analyser/pkg/httpClient"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StartProductsScheduler(ctx context.Context, dhm *mongodb.DBHandlerMongo, service *httpClient.Service) {
	startTime := time.Now()
	finish := make(chan struct{})

	log.Println("products scheduler: start products scheduler at", startTime)

	// query for products by page
	var currentMaxPage int
	queue := make(chan int, 100)

	// run first 5 pages
	for i := 1; i <= 99; i++ {
		currentMaxPage = i
		queue <- i
	}

	// listen and run next pages
	go func() {
		for {
			select {
			case page := <-queue:
				go func() {
					log.Println("products scheduler: start page", page)

					pageFinish := make(chan struct{})
					go func() {
						for {
							select {
							case <-pageFinish:
								// add new page to queue
								currentMaxPage++
								queue <- currentMaxPage
								return
							}
						}
					}()

					productsResponse, err := service.SendJSONRequest(ctx, http.MethodGet, "https://kaspi.kz/yml/product-view/pl/results?page="+strconv.Itoa(page), nil)
					if err != nil {
						log.Println("products scheduler: error while sending request for page", page, err)
						pageFinish <- struct{}{}
						return
					}

					data, ok := productsResponse["data"].([]interface{})
					if !ok {
						log.Println("products scheduler: error while parsing products response")
						pageFinish <- struct{}{}
						return
					}

					// if data array is empty, then we reached the end of the list
					if len(data) == 0 {
						finish <- struct{}{}
						pageFinish <- struct{}{}
						return
					}

					// parse products and save them to db with they offers
					for _, productData := range data {
						product, ok := productData.(map[string]interface{})
						if !ok {
							log.Println("products scheduler: error while parsing product")
							continue
						}

						id, ok := product["id"].(string)
						if !ok {
							log.Println("products scheduler: error while getting product id")
							continue
						}
						delete(product, "id")

						// save product to db
						if err = dhm.SaveProduct(ctx, id, product); err != nil && err != mongo.ErrNoDocuments {
							log.Println("products scheduler: error while saving product:", err)
							continue
						}

						// get offers for product
						offersResponse, err := service.SendJSONRequest(ctx, http.MethodPost, "https://kaspi.kz/yml/offer-view/offers/"+id, map[string]interface{}{
							"cityId": "750000000",
							"limit":  64,
						})
						if err != nil {
							log.Println("products scheduler: error while sending request for page", page, err)
							continue
						}

						offers, ok := offersResponse["offers"].([]interface{})
						if !ok {
							log.Println("products scheduler: error while parsing offers response", offersResponse)
							continue
						}

						// save offers to db
						for _, offerData := range offers {
							offer, ok := offerData.(map[string]interface{})
							if !ok {
								log.Println("products scheduler: error while parsing offer")
								continue
							}

							merchantId, ok := offer["merchantId"].(string)
							if !ok {
								log.Println("products scheduler: error while getting offer merchant id")
								continue
							}

							if err = dhm.SaveOffer(ctx, id+"_"+merchantId, offer); err != nil && err != mongo.ErrNoDocuments {
								log.Println("products scheduler: error while saving offer:", err)
								continue
							}
						}
					}

					pageFinish <- struct{}{}
				}()
			}
		}
	}()

	select {
	case <-finish:
		log.Println("products scheduler: products scheduler is done, time:", time.Since(startTime))
		StartProductsScheduler(ctx, dhm, service)
	}
}
