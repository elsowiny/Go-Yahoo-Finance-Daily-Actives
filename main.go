package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

func main() {

	type stockInfo struct {
		Symbol        string
		Name          string
		Price         string
		Change        string
		ChangePercent string
		Volume        string
		AvgVolume     string
		MarketCap     string
		PE            string
		//FiftyTwoWeekRange string
	}
	headers := []string{}
	//slice of slice of strings
	allStocksSlice := [][]string{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36"),
		colly.AllowedDomains("finance.yahoo.com"),
		colly.MaxBodySize(0),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)

	// Set max Parallelism and introduce a Random Delay
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       500 * time.Millisecond,
	})

	log.Println("User Agent: ", c.UserAgent)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())

	})

	c.OnHTML("thead", func(e *colly.HTMLElement) {
		log.Println("Found <thead> element")
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			el.ForEach("th", func(_ int, el *colly.HTMLElement) {
				//fmt.Println(el.Text)
				headers = append(headers, el.Text)
			})
		})
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		//	log.Println("Found <tbody> element")
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			dataSlice := []string{}
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {
				//fmt.Println(el.Text)
				//if its empty dont add to slice
				dataSlice = append(dataSlice, el.Text)
			})

			//add to overall slice
			allStocksSlice = append(allStocksSlice, dataSlice)
		})
	})

	c.Visit("https://finance.yahoo.com/most-active/")

	c.Wait()

	fmt.Println(headers)
	//loop allStocksSlice
	for i := 0; i < len(allStocksSlice); i++ {
		//fmt.Println("Stock: ", allStocksSlice[i])
		for j := 0; j < len(allStocksSlice[i]); j++ {
			//map headers to values
			fmt.Println(headers[j], ": ", allStocksSlice[i][j])

		}
	}

}
