package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	gomail "gopkg.in/mail.v2"
)

func main() {
	headers := []string{}
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
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())

	})
	c.OnHTML("thead", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			el.ForEach("th", func(_ int, el *colly.HTMLElement) {
				//ignore 52 Week Range
				if el.Text != "52 Week Range" {
					headers = append(headers, el.Text)
				}
			})
		})
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			dataSlice := []string{}
			el.ForEach("td", func(_ int, el *colly.HTMLElement) {
				dataSlice = append(dataSlice, el.Text)
			})
			//remove the last element which is the empty string
			dataSlice = dataSlice[:len(dataSlice)-1]

			//add to overall slice
			allStocksSlice = append(allStocksSlice, dataSlice)
		})
	})

	c.Visit("https://finance.yahoo.com/most-active/")

	c.Wait()

	printData(allStocksSlice, headers)
	saveDataToFile(allStocksSlice, headers)

	//ask user if they want to send email after checking if a config file exists
	if _, err := os.Stat("config.email.env"); os.IsNotExist(err) {
		noFileFound()
	} else {
		fmt.Println("Do you want to send email? (y/n)")
		var input string
		fmt.Scanln(&input)
		if input == "y" {
			sendToEmail("Daily_Actives.csv")
		} else {
			fmt.Println("Email not sent")
		}
	}

}

// helper functions

func saveDataToFile(data [][]string, headers []string) {
	fmt.Println("Saving data to file... Removing old file if exists")
	//delete file if exists
	os.Remove("Daily_Actives.csv")
	//open file
	file, err := os.Create("Daily_Actives.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//write to file
	writer := csv.NewWriter(file)

	//write headers
	writer.Write(headers)

	writer.Flush()
	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush()

	fmt.Println("Data saved to file")
}

func printData(data [][]string, headers []string) {

	//print headers
	fmt.Println("Printing")
	fmt.Println(headers)
	for _, value := range data {
		fmt.Println(value)
	}
}

func LoadConfig(path string) (config EmailConfig, err error) {
	viper.SetConfigName("config.email")
	viper.AddConfigPath(path)
	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, err

}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func sendToEmail(fileName string) {
	//send email
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	//validate email
	if !validateEmail(config.EMAIL) {
		fmt.Println("Invalid email not sending")
		return
	}

	//create email
	m := gomail.NewMessage()
	m.SetHeader("From", config.EMAIL)
	m.SetHeader("To", config.EMAIL)
	m.SetHeader("Subject", "Daily Actives")
	m.SetBody("text/html", "Daily Actives")
	m.Attach(fileName)

	//convert string to int for SMTP_PORT
	port, err := strconv.Atoi(config.SMTP_PORT)
	if err != nil {
		log.Fatal(err)
		return
	}
	//create dialer
	d := gomail.NewDialer(config.SMTP_HOST, port, config.EMAIL, config.PASSWORD)

	//send email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Println("Email sent")

}

func noFileFound() {
	fmt.Println("No config file found...")
	fmt.Println("if you want to send an email..")
	//wait 2 seconds
	time.Sleep(1 * time.Second)
	fmt.Println("Please create a config file named config.email.env with the following format..")
	time.Sleep(1 * time.Second)
	fmt.Println("EMAIL=EMAIL_ADDRESS")
	fmt.Println("PASSWORD=PASSWORD")
	fmt.Println("SMTP_HOST=SMTP_HOST")
	fmt.Println("SMTP_PORT=SMTP_PORT")
}
