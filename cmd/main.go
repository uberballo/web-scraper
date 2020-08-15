package main

import (
	"fmt"
	"log"
	"os"

	"github.com/uberballo/web-scraper/cmd/scraper"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	url := os.Getenv("MAIN_URL")
	res := scraper.Scrape(url)
	fmt.Println(string(res))
}
