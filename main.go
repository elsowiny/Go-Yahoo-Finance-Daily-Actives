package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

func main() {

	headers := []string{}
	stockData := []string{}

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
		log.Println("Found <tbody> element")
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			//fmt.Println(el.Text, "test tr")
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {
				//log.Println(el.Text)
				stockData = append(stockData, el.Text)
			})
		})
	})

	c.Visit("https://finance.yahoo.com/most-active/")

	c.Wait()

	fmt.Println(headers)
	//fmt.Println(stockData)
	//parse the stockData into a map

	//loop headers

	for i := 0; i < len(headers); i++ {
		fmt.Println(headers[i], ":", stockData[i])
	}

}
