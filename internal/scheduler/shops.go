package scheduler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"kaspi-analyser/internal/mongodb"
	"log"
	"net/http"
	"time"
)

func StartShopsScheduler(ctx context.Context, dhm *mongodb.DBHandlerMongo) {
	startTime := time.Now()
	finish := make(chan struct{})

	go func() {
		log.Println("shops scheduler: start shop scheduler at", startTime)

		merchants, err := dhm.GetMerchantsFromOffers(ctx)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("shops scheduler: error while getting merchants from offers:", err)
			finish <- struct{}{}
			return
		}

		for i, merchant := range merchants {
			resp, err := http.Get("https://kaspi.kz/shop/info/merchant/" + merchant.ID + "/address-tab/")
			if err != nil {
				log.Println("shops scheduler: error while getting shops:", err)
				finish <- struct{}{}
				return
			}

			if resp.StatusCode != 200 {
				log.Printf("shops scheduler: status code error: %d %s\n", resp.StatusCode, resp.Status)
				finish <- struct{}{}
				return
			}

			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Printf("shops scheduler: error while parsing html: %s", err)
				finish <- struct{}{}
				return
			}
			resp.Body.Close()

			doc.Find(".merchant-profile__contact-text").Each(func(i int, s *goquery.Selection) {
				phone := s.Text()
				merchant.PhoneNumber = phone
			})
			doc.Find(".merchant-profile__data-create").Each(func(i int, s *goquery.Selection) {
				created := s.Text()
				created = created[19 : len(created)-3]
				merchant.CreatedAt = created
			})

			merchants[i] = merchant
		}

		// merchants to []interface
		shops := make([]interface{}, len(merchants))
		for i, merchant := range merchants {
			shops[i] = merchant
		}

		// save shops to db
		if err = dhm.SaveShops(ctx, shops); err != nil && err != mongo.ErrNoDocuments && !mongo.IsDuplicateKeyError(err) {
			log.Println("shops scheduler: error while saving shop:", err)
			finish <- struct{}{}
			return
		}

		finish <- struct{}{}
	}()

	for {
		select {
		case <-finish:
			log.Println("shops scheduler: finish shop scheduler in", time.Since(startTime))

			// sleep 3 min before start new shops scheduler
			time.Sleep(3 * time.Minute)
			StartShopsScheduler(ctx, dhm)

			return
		}
	}
}
