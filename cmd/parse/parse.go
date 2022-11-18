package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	f, err := excelize.OpenFile("C:\\Users\\Islam\\go\\src\\kaspi-analyser\\cmd\\parse\\kaspi shops.xlsx")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.SetCellValue("offers", "E1", "phone numbers")
	f.SetCellValue("offers", "F1", "created at")

	done := make(chan bool)

	go func() {
		select {
		case <-done:
			if err = f.Save(); err != nil {
				panic(err)
			}
			panic("done")
		}
	}()

	i := 2
	for {
		var wg sync.WaitGroup

		i2 := i
		for j := i2; j < i2+10; j++ {
			wg.Add(1)

			go func(j int) {
				defer wg.Done()

				select {
				case <-done:
					return
				default:
					log.Printf("processing %d", j)

					cell, err := f.GetCellValue("offers", "A"+strconv.Itoa(j))
					if err != nil {
						log.Printf("error while getting cell value: %v\n", err)
					}
					if cell == "" {
						done <- true
						return
					}

					response, err := http.Get("https://kaspi.kz/shop/info/merchant/" + cell + "/reviews-tab/?redirect=false")
					if err != nil {
						log.Printf("error while getting %s: %v\n", cell, err)
						return
					}

					doc, err := goquery.NewDocumentFromReader(response.Body)
					if err != nil {
						log.Printf("error while parsing %s: %v\n", cell, err)
						return
					}

					doc.Find("span.merchant-profile__contact-text").Each(func(i int, s *goquery.Selection) {
						f.SetCellValue("offers", "E"+strconv.Itoa(j), s.Text())
					})

					doc.Find("div.merchant-profile__data-create").Each(func(i int, s *goquery.Selection) {
						f.SetCellValue("offers", "F"+strconv.Itoa(j), s.Text()[30:len(s.Text())-4])
					})
				}
			}(j)

			i++
		}

		wg.Wait()
	}
}
