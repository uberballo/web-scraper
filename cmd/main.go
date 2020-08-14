package main

import (
	"fmt"
	"os"

	"github.com/uberballo/web-scraper/cmd/scraper"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	url := os.Getenv("MAIN_URL")
	res := scraper.Scrape(url)
	fmt.Println(res)
}
