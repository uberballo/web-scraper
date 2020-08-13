package scraper

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var testSuffix string = "porssi/porssikurssit/osake/NOKIA/tilinpaatos"
var baseUrl string = "https://www.kauppalehti.fi"

func Scrape(url string) {
	visible(url)
}

func scrapeData(url, suffix string) {

}

func visible(url string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res []string
	var nodes []*cdp.Node

	fmt.Println("Starting scraping")
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.OMXH-list`),
		chromedp.Nodes(`div[class="list-striped mb-1"] > div > a`, &nodes, chromedp.ByQueryAll),
	)

	fmt.Println("Done scraping")
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range nodes {
		res = append(res, i.AttributeValue("href")+"/tilinpaatos")
	}
	for _, i := range res {
		fmt.Println(i)
	}
}
