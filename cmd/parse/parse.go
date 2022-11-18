package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strconv"
)

func main() {
	f, err := excelize.OpenFile("C:\\Users\\Islam\\go\\src\\kaspi-analyser\\cmd\\parse\\kaspi shops.xlsx")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.SetCellValue("offers", "E1", "phone numbers")
	f.SetCellValue("offers", "F1", "created at")

	i := 2
	for {
		cell, err := f.GetCellValue("offers", "A"+strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		if cell == "" {
			break
		}

		log.Printf("processing %d", i)

		response, err := http.Get("https://kaspi.kz/shop/info/merchant/" + cell + "/reviews-tab/?redirect=false")
		if err != nil {
			log.Printf("error while getting %s: %v\n", cell, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Printf("error while parsing %s: %v\n", cell, err)
			continue
		}

		doc.Find("span.merchant-profile__contact-text").Each(func(i int, s *goquery.Selection) {
			f.SetCellValue("offers", "E"+strconv.Itoa(i), s.Text())
		})

		doc.Find("div.merchant-profile__data-create").Each(func(i int, s *goquery.Selection) {
			f.SetCellValue("offers", "F"+strconv.Itoa(i), s.Text()[18:len(s.Text())-3])
		})

		i++
	}

	if err = f.Save(); err != nil {
		panic(err)
	}
}
