package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var testSuffix string = "porssi/porssikurssit/osake/NOKIA/tilinpaatos"
var baseUrl string = "https://www.kauppalehti.fi"
var url string = "https://www.kauppalehti.fi/porssi/kurssit/XHEL"

func Scrape(url string) {
	scrapeData(url, "")
}

func scrapeData(url, suffix string) {
	collectElements(url, `.OMXH-list`, "el")
}

func visible(url, element string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	timeoutContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := chromedp.Run(timeoutContext,
		chromedp.Navigate(url),
		chromedp.WaitVisible(element),
	)

	if err == nil {
		timeoutContext.Done()
		return nil
	}

	return err
}

func collectElements(url, containerElement, element string) {

	var res []string
	var nodes []*cdp.Node
	var err error

	fmt.Println("Starting scraping")

	err = visible(url, containerElement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("container element found")

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	fmt.Println("new cont")

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
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
